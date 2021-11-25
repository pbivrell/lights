import '../App.css';
import Button from "react-bootstrap/Button";
import React, { useState } from "react";
import LightsDetect from "../components/LightDetect";
import LightsIP from "../components/LightIP";
import { RgbColorPicker } from "react-colorful";
import axios from "axios";
import Spinner from "react-bootstrap/Spinner";

function Light () {

  	const [color, setColor] = useState({ r: 200, g: 150, b: 35 });
	const [ip, setIP] = useState("192.168.86.126");

	var last = Date.now();
	function onColorChange(color) {

		let now = Date.now();
		console.log(color);
		if (now - last < 500) {
			console.log("ignoring time");
			return;
		}
		last = now;
		setColor(color);
		axios.get(`http://${ip}/color?r=${color.r}&g=${color.g}&b=${color.b}`).catch((resp) => {
			console.log("TODO")
                });
	}


	return (
	<div className="App">
		<header className="App-header">
			<RgbColorPicker color={color} onChange={onColorChange} />
		</header>
	</div>
  	);
}


export default Light;
