package ssh

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client wraps golang.org/x/crypto/ssh for Proxmox node SSH access.
type Client struct {
	host    string
	port    int
	user    string
	sshConn *ssh.Client
}

// AuthMethod describes how to authenticate.
type AuthMethod string

const (
	AuthPassword AuthMethod = "password"
	AuthKey      AuthMethod = "key"
)

// Config holds SSH connection parameters.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	KeyPath  string
	Method   AuthMethod
}

// NewClient creates and connects an SSH client.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Port == 0 {
		cfg.Port = 22
	}
	if cfg.Method == "" {
		cfg.Method = AuthPassword
	}

	var authMethods []ssh.AuthMethod
	switch cfg.Method {
	case AuthPassword:
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	case AuthKey:
		keyData, err := os.ReadFile(cfg.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("reading SSH key: %w", err)
		}
		signer, err := ssh.ParsePrivateKey(keyData)
		if err != nil {
			return nil, fmt.Errorf("parsing SSH key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	sshCfg := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := ssh.Dial("tcp", addr, sshCfg)
	if err != nil {
		return nil, fmt.Errorf("SSH dial %s: %w", addr, err)
	}

	return &Client{
		host:    cfg.Host,
		port:    cfg.Port,
		user:    cfg.User,
		sshConn: conn,
	}, nil
}

// Run executes a command and returns combined stdout/stderr.
func (c *Client) Run(cmd string) (string, error) {
	session, err := c.sshConn.NewSession()
	if err != nil {
		return "", fmt.Errorf("creating SSH session: %w", err)
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	session.Stderr = &buf

	if err := session.Run(cmd); err != nil {
		return buf.String(), fmt.Errorf("running command %q: %w", cmd, err)
	}
	return buf.String(), nil
}

// Close closes the underlying SSH connection.
func (c *Client) Close() error {
	return c.sshConn.Close()
}

// IsAlive checks if the SSH connection is still alive with a keepalive.
func (c *Client) IsAlive() bool {
	_, _, err := c.sshConn.SendRequest("keepalive@openssh.com", true, nil)
	return err == nil
}

// Dial is a helper to check SSH connectivity without creating a full client.
func Dial(host string, port int) error {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
