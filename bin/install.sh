echo "Installing OWBot on system..."
cp owbot.service /etc/init.d/

echo "Enabling service on router start..."
/etc/init.d/owbot.service enable

/etc/init.d/owbot.service start &

echo "Your OWbot is installed and started."

