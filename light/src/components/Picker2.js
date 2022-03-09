import './Picker2.css';
import React, { useState, useEffect } from "react";
import { HuePicker, AlphaPicker } from 'react-color'
import bulb from '../res/Bulb_crop.png';
import offbulb from '../res/OffBulb_crop.png';
import disconnectedbulb from '../res/DisconnectedBulb_crop.png';
import axios from "axios";
import {displayError} from "../lib/hooks";

function Picker2({name, light, initalState}) {

	useEffect(() => {
		getLight();
	}, [light])


	const [bulbState, setBulbState] = useState(initalState);

  	const [color, setColor] = useState({ hsl: {a: 1}, rgb: {a: 1}, hex: ""});
	
	function getLight() {
		if (!light || !light.IP) {
			return
		}
		axios.get(`http://${light.IP}/status`, {timeout: 1000})
		.then((resp) =>{
			setBulbState(resp.status ? 0 : 1);
		})
		.catch((e) => {
			setBulbState(2);
			displayError(e);
		});
	}

	function pushLight(pushColor) {
		let data = generateBinary(pushColor, light.Count);
                let blob = new Blob([data], {type: "application/octet-stream"});
                let bodyFormData = new FormData();
                bodyFormData.append("file", blob, `${pushColor.r}-${pushColor.g}-${pushColor.b}-${pushColor.a}.bin`);

                axios({
                	method: "post",
                	url: `http://${light.IP}/upload`,
                	data: bodyFormData,
                	headers: { "Content-Type": "multipart/form-data" },
                })
                .then(function (response) {
                	console.log(response);
                })
                .catch(function (response) {
                	displayError(response)
                });
	}

	function pushToggle() {
		
		axios.get(`http://${light.IP}/toggle`, {timeout: 1000})
		.then((resp) => {
			console.log("toggled");
		})
		.catch((e) => {
        		displayError(e);
		});
	}

	function toggleBulb() {
		if (bulbState == 2) {
			getLight();
			return;
		}else if(bulbState == 1) {
			setBulbState(0);
		}else {
			setBulbState(1);
		}

		pushToggle();
	}

	function onColorChange(hue, c) {

		console.log(c);

		let hsl = c.hsl;
		let rgb = c.rgb;

		if(!hue) {
			hsl.a = color.hsl.a;
			rgb.a = color.rgb.a;
		}
		setColor({hsl: hsl, rgb: c.rgb, hex: `rgba(${c.rgb.r},${c.rgb.g},${c.rgb.b},${hsl.a}`});
		
		pushLight(c.rgb);
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
			data[(i + 1) * 8 + 1] = parseInt(color.r * color.a) ;
			data[(i + 1) * 8 + 2] = parseInt(color.g * color.a);
			data[(i + 1) * 8 + 3] = parseInt(color.b * color.a);
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

	function getToolTip() {
		if (!light || Object.keys(light).length === 0) {
			return ""
		}
		return `ID: ${light.ID}
IP: ${light.IP}
Version: ${light.Version}
${new Date(light.Updated).toLocaleDateString() +" " +  new Date(light.Updated).toLocaleTimeString()}`
	}

	return (
	<div className="Picker">

		<div className="ButtonText" data-tooltip={getToolTip()}>{name.alias}</div>
		<div className="Button" style={{background: color.hex}} onClick={toggleBulb}>
			<img className="ButtonImage" src={bulbState === 0 ? bulb : bulbState === 1 ? offbulb : disconnectedbulb }></img>
		</div>
		{ bulbState === 0 ? 
			<>
				<br/>
				<HuePicker className="ButtonImage" color={color.hsl} onChangeComplete={(c) => onColorChange(false, c)}/>
				<br/>
                		<AlphaPicker className="ButtonImage" color={color.hsl} onChangeComplete={(c) => onColorChange(true, c)}/>
			</>
		: 
			<></>
		}
	</div>
  	);
}


export default Picker2;
