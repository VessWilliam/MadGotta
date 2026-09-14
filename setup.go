package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// asset describes a single file to fetch and where it belongs.
type asset struct {
	name string
	path string
	url  string
	mode os.FileMode
}

func main() {
	fmt.Println("========================================")
	fmt.Println("Setup (Tailwind, HTMX, Alpine)")
	fmt.Println("========================================")

	root := projectRoot()
	staticDir := filepath.Join(root, "web", "static")
	jsDir := filepath.Join(staticDir, "js")
	cssDir := filepath.Join(staticDir, "css")
	inputCSS := filepath.Join(cssDir, "input.css")
	outputCSS := filepath.Join(cssDir, "output.css")
	tailwindBin := filepath.Join(staticDir, tailwindBinaryName())

	assets := []asset{
		{
			name: "Tailwind CLI",
			path: filepath.Join(staticDir, tailwindBinaryName()),
			url:  tailwindDownloadURL(),
			mode: 0755,
		},
		{
			name: "htmx",
			path: filepath.Join(jsDir, "htmx.min.js"),
			url:  "https://unpkg.com/htmx.org@latest/dist/htmx.min.js",
			mode: 0644,
		},
		{
			name: "Alpine.js",
			path: filepath.Join(jsDir, "alpine.min.js"),
			url:  "https://unpkg.com/alpinejs@latest/dist/cdn.min.js",
			mode: 0644,
		},
	}

	for _, a := range assets {
		if err := a.ensure(); err != nil {
			fmt.Fprintf(os.Stderr, "error setting up %s: %v\n", a.name, err)
			os.Exit(1)
		}
	}

	if err := ensureInputCSS(inputCSS); err != nil {
		fmt.Fprintf(os.Stderr, "error creating input.css: %v\n", err)
		os.Exit(1)
	}

	if err := buildTailwind(tailwindBin, inputCSS, outputCSS); err != nil {
		fmt.Fprintf(os.Stderr, "error building CSS: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAll assets are downloaded and ready.")
}

func (a asset) ensure() error {
	if _, err := os.Stat(a.path); err == nil {
		fmt.Printf("[%s] Found: %s\n", a.name, a.path)
		return nil
	}

	fmt.Printf("[%s] Downloading from:\n  %s\n", a.name, a.url)

	if err := os.MkdirAll(filepath.Dir(a.path), 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if err := download(a.url, a.path, a.mode); err != nil {
		return err
	}

	fmt.Printf("[%s] Saved to %s\n", a.name, a.path)
	return nil
}

func download(url, path string, mode os.FileMode) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	tmp := path + ".tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		os.Remove(tmp)
		return fmt.Errorf("write file: %w", err)
	}
	out.Close()

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("finalize file: %w", err)
	}

	return nil
}

func projectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	if filepath.Base(cwd) == "config" {
		return filepath.Dir(cwd)
	}
	return cwd
}

func tailwindBinaryName() string {
	if runtime.GOOS == "windows" {
		return "tailwindcss.exe"
	}
	return "tailwindcss"
}

func tailwindDownloadURL() string {
	const base = "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/"

	switch runtime.GOOS {
	case "windows":
		return base + "tailwindcss-windows-x64.exe"
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return base + "tailwindcss-macos-arm64"
		}
		return base + "tailwindcss-macos-x64"
	default:
		if runtime.GOARCH == "arm64" {
			return base + "tailwindcss-linux-arm64"
		}
		return base + "tailwindcss-linux-x64"
	}
}

func buildTailwind(bin, input, output string) error {
	fmt.Printf("[Tailwind] Building %s -> %s\n", input, output)

	cmd := exec.Command(bin, "-i", input, "-o", output, "--minify")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tailwind build: %w", err)
	}

	fmt.Printf("[Tailwind] Saved to %s\n", output)
	return nil
}
func ensureInputCSS(path string) error {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("[input.css] Found: %s\n", path)
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	content := "@import \"tailwindcss\";\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write input.css: %w", err)
	}

	fmt.Printf("[input.css] Created: %s\n", path)
	return nil
}
