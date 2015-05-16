var deviceHostname = "localhost:8080"
var remoteHostname = "localhost:8080"

function uploadFiles(files, url)
{
    $.ajax({
        url: url,
        type: 'POST',
        data: JSON.stringify(files),
        cache: false,
        contentType: 'application/json; charset=UTF-8',
      });
}

function pollAndPush() {
	// Get from remote upload to device
    $.getJSON( "http://" + remoteHostname + "getfiles/callback=?", function( files ) {
    					console.log(files);
                        if(files.length > 0) {
                            uploadFiles(files, "http://" + deviceHostname + "/upload");
                        }
	});
	// vice versa
    $.getJSON( "http://" + deviceHostname + "/getfiles/callback=?", function( files ) {
        					console.log(files);
                            if(files.length > 0) {
                                uploadFiles(files, "http://" + remoteHostname + "/upload");
                            }
	});
}

(function () {
 
    // dynamicly loads javascript
    function loadScript(url, callback) {
 
        var script = document.createElement("script")
        script.type = "text/javascript";
 
        if (script.readyState) { //IE
            script.onreadystatechange = function () {
                if (script.readyState == "loaded" || script.readyState == "complete") {
                    script.onreadystatechange = null;
                    callback();
                }
            };
        } else { //Others
            script.onload = function () {
                callback();
            };
        }
 
        script.src = url;
        document.getElementsByTagName("head")[0].appendChild(script);
    }
 	
    loadScript("http://" + deviceHostname + "/jquery-2.1.4.js", function () {
        setInterval(pollAndPush, 2000);
	});
 
})();

