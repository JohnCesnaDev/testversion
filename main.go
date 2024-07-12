package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
)

const (
    githubRepo        = "utilisateur/repo"    // Remplacez par le repo GitHub
    versionFilePath   = "/chemin/vers/clone/version.txt" // Chemin vers le fichier de version dans le dépôt cloné
    cloneDir          = "/chemin/vers/clone/"       // Répertoire où le dépôt sera cloné
)

func getCurrentVersion() (string, error) {
    data, err := ioutil.ReadFile(versionFilePath)
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(data)), nil
}

func getLatestVersion() (string, error) {
    cmd := exec.Command("git", "ls-remote", "--tags", fmt.Sprintf("https://github.com/%s.git", githubRepo))
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    lines := strings.Split(string(output), "\n")
    if len(lines) == 0 {
        return "", fmt.Errorf("Aucune version trouvée")
    }

    lastLine := lines[len(lines)-2]
    parts := strings.Split(lastLine, "/")
    if len(parts) == 0 {
        return "", fmt.Errorf("Erreur lors de la récupération de la version")
    }

    return parts[len(parts)-1], nil
}

func cloneOrUpdateRepo() error {
    if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
        cmd := exec.Command("git", "clone", fmt.Sprintf("https://github.com/%s.git", githubRepo), cloneDir)
        if err := cmd.Run(); err != nil {
            return err
        }
    } else {
        cmd := exec.Command("git", "-C", cloneDir, "pull")
        if err := cmd.Run(); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    currentVersion, err := getCurrentVersion()
    if err != nil {
        if !os.IsNotExist(err) {
            fmt.Println("Erreur lors de la lecture de la version actuelle:", err)
            return
        }
    }

    latestVersion, err := getLatestVersion()
    if err != nil {
        fmt.Println("Erreur lors de la récupération de la dernière version:", err)
        return
    }

    if currentVersion != latestVersion {
        fmt.Printf("Nouvelle version trouvée : %s. Mise à jour...\n", latestVersion)
        if err := cloneOrUpdateRepo(); err != nil {
            fmt.Println("Erreur lors de la mise à jour du dépôt:", err)
            return
        }

        fmt.Println("Le dépôt a été mis à jour.")
    } else {
        fmt.Println("Votre application est déjà à jour.")
    }
}
