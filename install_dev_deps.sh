echo "Installing Go development dependencies..."
go install github.com/pressly/goose/v3/cmd/goose@latest;
go install github.com/air-verse/air@latest;
echo "Go development dependencies installed successfully."
echo "Adding bash alias"
echo "alias ll=\"ls -l\"" >> ~/.bashrc
echo "alias la=\"ls -la\"" >> ~/.bashrc
echo "alias cgoose=\"./build/goose\"" >> ~/.bashrc