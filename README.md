<h1 align="center">pre-requisites: make sure you have go installed:</h2>
<h2 align="center">
    <code> go version </code>
</h2>
<h3 align="center">(requires go version 1.25 or higher.)</h3>

<br>
<h1 align="center">getting started:</h2>

<h2 align="center">setting headers while testing</h2>

<h3 align="center"> 
    <code> go run . [-H key:value] [URL] </code>
</h>


<br>
<h2 align="center">testing with different http methods</h2>
<h3 align="center">
    <code> go run . [-x METHOD] [URL] </code>
</h3>

<br>

<h2 align="center">with streaming enabled:</h2>

<h3 align="center">
    <code> go run . [-H key:value] [-stream true/false] [URL] </code>
</h3>

<h3 align="center">(you can include 1 of many flags while testing aswell, but make sure to always include the URL at the end.)</h3>

<br>

<h2 align="center">to allow response body logging:</h2>

<h3 align="center">
    <code> go run . [-b true/false -s 1024] [URL] </code>
</h3>

<h4 align="center">(you can set a limit on the response body size to log, as with the -s flag, if you want it to be the default (1 KiB), simply do not include the flag and only include the -b flag).</h4>
<br>
<h2 align="center">concurrent testing</h2>
<h3 align="center">(sending numerous requests simultaneously):</h4>
<h3 align="center">
    <code> go run . [-c (number)] [URL] </code>
</h3>
<h4 align="center">be sure to include a reasonable number of requests to be sent.</h4>
<br>
<h2 align="center">session mode:</h2>

<h3 align="center">Instead of manually pasting, you can enter session mode to dynamically store urls, headers, etc in a variable and then use that variable to test API's (in which the variable will hold the header/url etc).</h3>

<br>
<h3 align="center">to enter session mode:</h3>

<h3 align="center">
    <code> go run . [-session true/false] </code>
</h3>

<h3 align="center">when entering session mode, the url does not need to be present at the end.</h3>
<br>
<h2 align="center">session mode tutorial:</h2>
<h3 align="center">start off by setting a variable to hold a url/header or such.</h3>
<h3 align="center">
    <code> VAR [VarName] [Value]</code>
</h3>
<h3 align="center">and to retrieve a variable:</h3>
<h3 align="center">
    <code> GET [VarName] </code>
</h3>
<h3 align="center">and it returns the value of whatever the variable stores.</h3>
<h3 align="center">to delete a variable, run:</h3>
<h3 align="center">
    <code> DEL [VarName] </code>
</h3>
<h3 align="center">to see all stored variables:</h3>
<h3 align="center">
    <code> FETCH </code>
</h3>
<h3 align="center">view all available commands whilst your in session mode:</h3>
<h3 align="center">
    <code> HELP </code>
</h3>
<br>
<h3 align="center">using the variables you stored for API testing:</h3>

<h4 align="center">(for standard get requests):</h4>
<h3 align="center">
    <code> TEST [VarName] </code>
</h3>

<h4 align="center">if you a custom header aswell, do so in the following format:</h4>
<h3 align="center">
    <code> TEST [VarName] [-h headerName:value] </code>
</h3>

<h4 align="center">you can also configure the method as so:</h4>
<h3 align="center">
    <code> TEST [VarName] [-x GET] </code>
</h3>
<h4 align="center">and post request like so:</h4>
<h3 align="center">
    <code> TEST [VarName] [-x POST -D [bodyData]] </code>
</h3>
<h4 align="center">for post requests, when setting a request body you can do -F for form related data, or -D for normal json data as shown above, you can use the -F flag like so:</h4>
<h3 align="center">
    <code> TEST [VarName] [-x POST  -F [title:value]] </code>
</h3>
<h3 align="center">(compatible with any http method.)</h3>
<br>

<h2 align="center">sending multiple requests:</h2>

<h3 align="center">for post request:</h3>
<h3 align="center">
    <code> SPAM [VarName] [-x POST] [-c 10] [-d {"name": "alex"}] </code>
</h3>
<h4 align="center">be sure to include the "-c" flag as it is the number of requests that are to be sent to the api, with the json data that should be provided in the "-d" flag. </h4>

<h3 align="center">and for get requests:</h3>
<h3 align="center">
    <code> SPAM [VarName] [-x GET] [-c 10] </code>
</h3>
<h4 align="center">important to note that you can include an variable as the argument to the "-c" flag, though make sure the variables value is an integer.</h4>
<br>
<h2 align="center">to exit session mode:</h3>
<h3 align="center">
    <code> EXIT </code>
</h3>
<br>
