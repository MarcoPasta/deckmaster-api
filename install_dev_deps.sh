echo "Installing Go development dependencies..."
go install github.com/pressly/goose/v3/cmd/goose@latest;
go install github.com/air-verse/air@latest;
go install github.com/chromedp/chromedp@latest;
echo "Go development dependencies installed successfully."
echo "Installing Python development dependencies..."
echo "Adding bash alias"
echo "alias ll=\"ls -l\"" >> ~/.bashrc
echo "alias la=\"ls -la\"" >> ~/.bashrc
curl -LsSf https://astral.sh/uv/install.sh | sh