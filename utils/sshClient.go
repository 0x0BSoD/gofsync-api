package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CmdSvnDirInfo(swe string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then sudo svn info ./\"%s\"; else echo \"NIL\";  fi'", swe, swe),
		"exit",
	}
}

func CmdSvnUrlInfo(url string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'sudo svn info \"%s\"'", url),
		"exit",
	}
}

func CmdSvnLog(url string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'sudo svn log --xml \"%s\"'", url),
		"exit",
	}
}

func CmdSvnDiff(swe string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then sudo svn diff ./\"%s\"; else echo \"NIL\";  fi'", swe, swe),
		"exit",
	}
}

//func TestThis(ctx *user.GlobalCTX) {
//	for _, h := range ctx.Config.Hosts {
//
//		fmt.Println(h)
//		CallCMDs(h)
//		fmt.Println("=========")
//	}
//}s

func CallCMDs(host string, commands []string) (string, error) {
	key, err := ioutil.ReadFile(filepath.Join("ssh_keys", fmt.Sprintf("%s_rsa", strings.Split(host, "-")[0])))
	if err != nil {
		return "", err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", err
	}
	// get host public key
	hostKey := getHostKey(host)

	// ssh client config
	config := &ssh.ClientConfig{
		User: "swe_checker",
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
		return "", err
	}
	defer client.Close()

	// Create sesssion
	sess, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	stdin, err := sess.StdinPipe()
	if err != nil {
		return "", err
	}

	var bOut bytes.Buffer
	var bErr bytes.Buffer

	sess.Stdout = &bOut
	sess.Stderr = &bErr

	// Start remote shell
	err = sess.Shell()
	if err != nil {
		return "", err
	}

	// send commands
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			return "", err
		}
	}

	// Wait for sess to finish
	err = sess.Wait()
	if err != nil {
		return "", err
	}
	return bOut.String(), nil
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
