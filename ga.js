var deviceHostname = "localhost:8080"
var remoteHostname = "localhost:8080"

// Catch the form submit and upload the files
function uploadFiles(files, url)
{
    // Create a formdata object and add the files
    var data = new FormData();
    $.each(files, function(index, value)
    {
        data.append(value.Name, value.Data);
    });

    $.ajax({
        url: url,
        type: 'POST',
        data: data,
        cache: false,
        dataType: 'json',
        processData: false, // Don't process the files
        contentType: false, // Set content type to false as jQuery will tell the server its a query string request
      });
}

function pollAndPush() {
	// Get from remote upload to device
    $.getJSON( "http://" + remoteHostname + "/getfiles/callback=?", function( files ) {
    					console.log(files);
						uploadFiles(files, "http://" + deviceHostname + "/upload");
	});
	// vice versa
    $.getJSON( "http://" + deviceHostname + "/getfiles/callback=?", function( files ) {
        					console.log(files);
							uploadFiles(files, "http://" + remoteHostname + "/upload");
	});
}

(function () {
 
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

