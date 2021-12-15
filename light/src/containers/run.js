import '../App.css';
import React, { useState, useEffect } from "react";

function Run() {

	//const [ip, setIP] = useState("192.168.86.136");
	//const [ip, setIP] = useState("192.168.86.144");
	const [ip, setIP] = useState("192.168.86.152");

	return (
	<div className="App">
		<header className="App-header">
			<form method='POST' enctype='multipart/form-data' action={"http://" + ip + "/upload"}>
     				<input type="file" id="myFile" name="filename"></input>
     				<input type="submit"></input>
   			</form>
		</header>
	</div>
  	);
}


export default Run;
