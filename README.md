# Gestionnaire de Presse-papiers

Une application de gestion de presse-papiers simple, construite avec Go et Fyne, permettant de conserver un historique des copies et de les recoller facilement.

## Prérequis

- **Go 1.19 ou version supérieure**
- **Git**

## Installation

### 1. Installer Go (si non déjà installé)

#### Sous Ubuntu :

```bash
sudo apt remove golang-go # Supprimez les anciennes versions si présentes
sudo apt autoremove

# Ajouter le PPA pour les versions récentes de Go
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go

# Vérifier la version de Go
go version
```

Vous devriez voir une sortie similaire à : `go version go1.19.x linux/amd64`

#### Sous Windows ou macOS :

- Rendez-vous sur [https://go.dev/dl/](https://go.dev/dl/)
- Téléchargez et installez la version la plus récente de Go.

### 2. Installer les dépendances

Pour l'utilisation de Fyne:

```bash
sudo apt install libgl1-mesa-dev xorg-dev
```


Cela installera automatiquement les dépendances comme `fyne.io/fyne/v2` et `github.com/atotto/clipboard`:

```bash
go mod tidy
```
### 3. Ajouter un raccourci clavier

#### i. installer xbindkeys:

```bash
sudo apt install xbindkeys
```

#### ii. Configurer le raccourci:

```bash
sudo nano ~/.xbindkeysrc
```

Ajouter:

```bash
"/Chemin/Vers/Le/Projet/clipboard_ui"
    Control+Mod2+Mod4 + v
```

#### iii. Activer la configuration:

```bash
xbindkeys
```

### 4. Exécuter l'application

```bash
go build -o clipboard_daemon.go
go build -o cliboard_ui.go
```

L'application s'ouvrira avec une interface graphique simple pour gérer votre historique de presse-papiers.

### 5. Compiler l'application (optionnel)

Pour créer un exécutable :

```bash
go build -o gestionnaire-presse-papiers
./gestionnaire-presse-papiers
```

### 6. Ajouter clipboard_daemon aux fichiers de démarrage

#### i.Tester le daemon:

```bash
./clipboard_daemon &
```

#### ii. Crée un fichier de service:

```bash
sudo nano /etc/systemd/system/clipboard_daemon.service
```

Ajouter le contenu suivant:

```bash
[Unit]
Description=Clipboard Daemon
After=network.target

[Service]
ExecStart=/chemin/vers/clipboard_daemon
Restart=always
User=%I
WorkingDirectory=/chemin/vers/votre/projet

[Install]
WantedBy=multi-user.target
```
Remplacer %I par votre nom d'utilisateur ainsi que les chemins vers le working directory et l'executable

#### iii. Activer le service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable clipboard_daemon
sudo systemctl start clipboard_daemon
```

#### iv. Verifier le fonctionnement:

```bash
sudo systemctl status clipboard_daemon
sudo reboot
ps aux | grep clipboard_daemon
```


## Fonctionnalités

- Surveillance automatique du presse-papiers.
- Historique des 20 dernières copies.
- Double-clic sur un élément pour le recoller rapidement.

## Problèmes courants

- **Problèmes de permissions sur Linux :**
Essayez d'exécuter l'application avec des permissions élevées si le presse-papiers ne réagit pas :

```bash
sudo go run main.go
```
