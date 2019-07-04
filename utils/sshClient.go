package utils

import (
	"bufio"
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func TestThis(ctx *user.GlobalCTX) {
	for _, h := range ctx.Config.Hosts {
		fmt.Println(h)
		CallCMD(h)
		fmt.Println("=========")
	}
}

func CallCMD(host string) {

	key, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	if err != nil {
		Error.Printf("unable to read private key: %v", err)
		return
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		Error.Printf("unable to parse private key: %v", err)
		return
	}
	// get host public key
	hostKey := getHostKey(host)

	// ssh client config
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// verify host public key
		HostKeyCallback: ssh.FixedHostKey(hostKey),
		Timeout:         5 * time.Second,
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		Error.Printf("unable to connect: %v", err)
		return
	}
	defer client.Close()

	// Create sesssion
	sess, err := client.NewSession()
	if err != nil {
		Error.Printf("Failed to create session: %s", err)
		return
	}
	defer sess.Close()

	stdin, err := sess.StdinPipe()
	if err != nil {
		Error.Println(err)
		return
	}

	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// Start remote shell
	err = sess.Shell()
	if err != nil {
		Error.Println(err)
		return
	}

	// send the commands
	commands := []string{
		"pwd",
		"whoami",
		"echo 'bye'",
		"exit",
	}
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			Error.Println(err)
			return
		}
	}

	// Wait for sess to finish
	err = sess.Wait()
	if err != nil {
		Error.Println(err)
		return
	}
}

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		Error.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				Error.Printf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		Error.Printf("no hostkey found for %s", host)
	}

	return hostKey
}
