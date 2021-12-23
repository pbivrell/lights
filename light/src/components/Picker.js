import '../App.css';
import Button from "react-bootstrap/Button";
import React, { useState } from "react";
import LightsDetect from "../components/LightDetect";
import LightsIP from "../components/LightIP";
import { RgbColorPicker } from "react-colorful";
import axios from "axios";
import Spinner from "react-bootstrap/Spinner";

function Picker({lightData, lights}) {

  	const [color, setColor] = useState({ r: 200, g: 150, b: 35 });
	//const [ip, setIP] = useState({ ip: "192.168.86.162", pixels: 50 });

	var last = Date.now();
	function onColorChange(color) {

		let now = Date.now();
		console.log(color);
		if (now - last < 100) {
			console.log("ignoring time");
			return;
		}
		last = now;
		setColor(color);
		
	
		console.log(lightData);
		console.log(lights);
		console.log(lights.lights);
		
		if (!lights || !lights.lights || !lightData) {
			return;
		}
		
		lights.lights.forEach((light) => {
			console.log("looping", light);
			let data = generateBinary(color, lightData.lights[light].pixels);
                	let blob = new Blob([data], {type: "application/octet-stream"});
			let bodyFormData = new FormData();
                	bodyFormData.append("file", blob, "dummy.bin");

                	axios({
                		method: "post",
                		url: `http://${lightData.lights[light].ip}/upload`,
                		data: bodyFormData,
                		headers: { "Content-Type": "multipart/form-data" },
                	})
                	.then(function (response) {
                		//handle success
                		console.log(response);
                	})
                	.catch(function (response) {
                		//handle error
                		console.log(response);
                	});
		});

		//axios.get(`http://${ip}/color?r=${color.r}&g=${color.g}&b=${color.b}`).catch((resp) => {
		//	console.log("TODO")
                //});
	}
	
	
	function generateBinary(color, pixels) {
		
		let data = new Uint8Array((2 + pixels) * 8);
		let bytes = getInt64Bytes(pixels);

		let count = (bytes[bytes.length-2] << 8) | bytes[bytes.length-1];
		
		data[0] = 0x3;
		data[1] = bytes[bytes.length-1];
		data[2] = bytes[bytes.length-2];

		for (let i = 0; i < pixels; i++) {
			let bytes = getInt64Bytes(i);
		 	data[(i + 1) * 8] = 0x1;
			data[(i + 1) * 8 + 1] = color.r;
			data[(i + 1) * 8 + 2] = color.g;
			data[(i + 1) * 8 + 3] = color.b;
			data[(i + 1) * 8 + 4] = bytes[bytes.length-1];
			data[(i + 1) * 8 + 5] = bytes[bytes.length-2];
		}

		data[8 + (pixels * 8)] = 0x2;
		data[8 + (pixels * 8) + 1] = 200;
		data[8 + (pixels * 8) + 2] = 0;
		data[8 + (pixels * 8) + 3] = 0;
		data[8 + (pixels * 8) + 4] = 0;


		console.log(data);
		return data;

	}

	function getInt64Bytes(x) {
  		let y= Math.floor(x/2**32);
  		return [y,(y<<8),(y<<16),(y<<24), x,(x<<8),(x<<16),(x<<24)].map(z=> z>>>24)
	}

	return (
	<div>
            <RgbColorPicker color={color} onChange={onColorChange} />
            <p>{ color.r}, {color.g}, {color.b}</p>
	</div>
  	);
}


export default Picker;
