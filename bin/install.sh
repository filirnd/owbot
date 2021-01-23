version=1.0.0
echo ""
echo "   \\     /                   #### OWbot ####  v. "$version
echo "   _\\___/_"
echo " /______ /|  Yet another telegram bot, but for your router."
echo "|_°_____|/   Made with <3 by Filirnd (https://github.com/filirnd/)"
echo ""
echo ""

echo "Installing bot on system..."
cp owbot.service /etc/init.d/

echo "Enabling service on router start..."
/etc/init.d/owbot.service enable

/etc/init.d/owbot.service start

echo "Your OWbot is installed and started."

