Vagrant.configure("2") do |config|
  config.vm.box = "win10"
  config.vm.provision "shell", path: "provision.ps1"
  config.vm.network "public_network", ip: "10.199.199.99"
  config.vm.network :forwarded_port, host: 3389, guest: 3389, id: "rdp", auto_correct: true
  config.vm.communicator = "winrm"
  config.winrm.password = "Passw0rd!"
  config.winrm.username = "IEUser"
end
