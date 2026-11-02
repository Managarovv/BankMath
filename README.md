****Start****    
**1.**   
git clone https://github.com/Managarovv/BankMath.git  
cd BankMath  
**2.**  
wget https://github.com/Managarovv/BankMath/releases/latest/download/BankMath-linux-amd64 or BankMath-linux-arm64 -O BankMath  
chmod +x BankMath  
**3.**  
sudo tee /etc/systemd/system/bankmath.service > /dev/null <<EOF  
[Unit]  
Description=BankMath web service  
After=network.target  

[Service]
Type=simple  
User=$USER  
WorkingDirectory=$(pwd)  
ExecStart=$(pwd)/BankMath  
Restart=always  
RestartSec=5  
  
[Install]  
WantedBy=multi-user.target  
EOF    
**4.**  
sudo systemctl daemon-reload  
sudo systemctl enable bankmath.service  
sudo systemctl start bankmath.service  
