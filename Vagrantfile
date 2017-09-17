Vagrant.configure("2") do |config|
	config.vm.box = "Ubuntu1604"
	config.vm.box_url = "https://www.dropbox.com/s/g5tzb35b58sr6tr/ubuntu1604lts5110.box?dl=1"
	config.ssh.insert_key = false	# Avoid that vagrant removes default insecure key

	# Host manager setup
	config.hostmanager.enabled				= true
	config.hostmanager.manage_host			= true
	config.hostmanager.manage_guest			= true
	config.hostmanager.ignore_private_ip	= false
	config.hostmanager.include_offline		= true

	config.vm.define 'SessionService' do |devbox|
        devbox.vm.network "private_network", ip: "192.168.2.10"
        devbox.vm.hostname = "session.dev"

        devbox.vm.provider "virtualbox" do |vb|
			vb.name = "SessionService"
			# Don't boot with headless mode
# 			vb.gui = true

			# Use VBoxManage to customize the VM. For example to change memory:
			vb.customize ["modifyvm", :id, "--memory", "1024"]
        end

        # Enable provisioning with chef solo, specifying a cookbooks path, roles
        # path, and data_bags path (all relative to this Vagrantfile), and adding
        # some recipes and/or roles.
        #
        devbox.vm.provision "chef_solo" do |chef|
            chef.cookbooks_path = "./vendor/rebel-l/sisa/cookbooks"
            chef.roles_path = "./vendor/rebel-l/sisa/roles"
            chef.environments_path = "./vendor/rebel-l/sisa/environments"
            chef.data_bags_path = "./vendor/rebel-l/sisa/data_bags"
            chef.add_role "DockerServer"
            chef.environment = "development"
            chef.add_recipe "GolangCompiler"

            # You may also specify custom JSON attributes:
            chef.json = {
                'Golang' => {
                    'version' => '1.9'
                },
                'System' => {
                    'Iptables' => {
                        'TCP' => {
                            'Ports' => [
                                4000
                            ]
                        }
                    }
                }
            }
        end
    end
end