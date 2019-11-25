package utils

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

// =====================================================================================================================
// COMMANDS
// =====================================================================================================================

// svn info --xml
func CmdSvnDirInfo(swe string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then sudo svn info --xml ./\"%s\"; else echo \"NIL\";  fi'", swe, swe),
		"exit $(echo $?)",
	}
}

// svn info
func CmdSvnUrlInfo(url string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'sudo svn info --xml \"%s\"'", url),
		"exit $(echo $?)",
	}
}

// svn log --xml
func CmdSvnLog(url string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'sudo svn log --xml \"%s\"'", url),
		"exit $(echo $?)",
	}
}

// svn update
func CmdSvnUpdate(name string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'sudo svn update \"%s\"'", name),
		fmt.Sprintf("bash -c 'sudo chown -R puppet:puppet %s'", name),
		fmt.Sprintf("bash -c 'sudo chmod -R 755 %s'", name),
		"exit $(echo $?)",
	}
}

// svn checkout
func CmdSvnCheckout(url, name string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'sudo svn checkout \"%s\"'", url),
		fmt.Sprintf("bash -c 'sudo chown -R puppet:puppet %s'", name),
		fmt.Sprintf("bash -c 'sudo chmod -R 755 %s'", name),
		"exit $(echo $?)",
	}
}

// svn diff
func CmdSvnDiff(swe string) []string {
	return []string{
		"cd /etc/puppet/environments",
		fmt.Sprintf("bash -c 'if [ -d \"./%s\" ]; then sudo svn diff ./\"%s\"; else echo \"NIL\";  fi'", swe, swe),
		"exit $(echo $?)",
	}
}

// =====================================================================================================================
// WRAPPER
// =====================================================================================================================
func CallCMDs(host string, commands []string) (string, error) {
	key, err := ioutil.ReadFile(filepath.Join("ssh_keys", fmt.Sprintf("%s_rsa", strings.Split(host, "-")[0])))
	if err != nil {
		return "", err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", err
	}

	// ssh client config
	config := &ssh.ClientConfig{
		User: "swe_checker",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// not verify host public key
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Create a session
	sess, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	// ###############################
	// ################

	stdin, err := sess.StdinPipe()
	if err != nil {
		//Warning.Println(err)
		return "", err
	}

	var bOut bytes.Buffer
	sess.Stdout = &bOut

	var bErr bytes.Buffer
	sess.Stderr = &bErr

	err = sess.Shell()
	if err != nil {
		Warning.Println("failed to start shell: ", err)
		return "", err
	}

	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			//Warning.Println(err)
			return "", err
		}
	}
	err = sess.Wait()
	if err != nil {
		//Warning.Println(err)
		return "", err
	}

	// ################
	// ###############################

	response := bOut.String()

	return response, nil
}
