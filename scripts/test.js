function run() {

    Process("plugins.psutil.cpu")
    Process("plugins.psutil.disk")
    Process("plugins.psutil.docker")
    Process("plugins.psutil.host")
    Process("plugins.psutil.load")
    Process("plugins.psutil.mem")
    Process("plugins.psutil.net")
    Process("plugins.psutil.process")
    Process("plugins.psutil.winservices")


}

