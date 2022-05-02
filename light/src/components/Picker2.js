import './Picker2.css';
import React, { useState, useEffect } from "react";
import { HuePicker, AlphaPicker } from 'react-color'
import bulb from '../res/Bulb_crop.png';
import offbulb from '../res/OffBulb_crop.png';
import disconnectedbulb from '../res/DisconnectedBulb_crop.png';
import axios from "axios";
import {displayError} from "../lib/hooks";
import LightSettings from '../containers/lightsettings.js';
import {ToggleButton, ToggleButtonGroup, DropdownButton, Dropdown, ListGroup} from "react-bootstrap";

function Picker2({name, hub}) {

	const [settings, setSettings] = useState(false);
	const [light, setLight] = useState({});
	const [patternMode, setPatternMode] = useState(true);
	const [patterns, setPatterns] = useState([]);
	const [selected, selectPattern] = useState("");

	useEffect(() => {
		getLight();
		getPatterns();
	}, [])


  	const [color, setColor] = useState({ hsl: {a: 1}, rgb: {a: 1}, hex: ""});
	
	function settingsCallback() {
		setSettings(false);
		getLight();
	}

	function getLight() {
		
		axios.get(`${process.env.REACT_APP_API_URL}/user/hub/${hub}/light/${name}`, {
			withCredentials: true,
		}).then((resp) => {
			setLight(resp.data);
			let pattern = resp.data.PatternID.split("~")[1].split(".")[0].split("-");
			console.log(pattern);
			setColor({
				rgb: {r: pattern[0], g: pattern[1], b: pattern[2], a: pattern[3]/100},
				hex: `rgba(${pattern[0]},${pattern[1]},${pattern[2]},${pattern[3]/100}`,
				hsl: {a: pattern[3]/100},

			});
		}).catch((e) => {
			displayError(e);
		});
	}

	function getPatterns() {
		setPatterns([
			{"alias": "spring", "id": 0, "creator": "demo"},
			{"alias": "spring", "id": 1, "creator": "demo"},
		]);
	}

	/*function pushLight(pushColor) {
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
	}*/

	function pushLight(pushColor) {
		axios.post(
			`${process.env.REACT_APP_API_URL}/user/hub/${hub}/light/${name}`, 
			{ PatternID: `~${pushColor.r}-${pushColor.g}-${pushColor.b}-${Math.floor(pushColor.a * 100)}.bin` }, 
			{withCredentials: true,}
		)
		.then((resp) => {
			console.log("toggled");
		})
		.catch((e) => {
        		displayError(e);
		});
	}

	function pushToggle(s) {
		axios.post(
			`${process.env.REACT_APP_API_URL}/user/hub/${hub}/light/${name}`,
			{Status: s}, 
			{withCredentials: true,}
		)
		.then((resp) => {
			console.log("toggled");
		})
		.catch((e) => {
        		displayError(e);
		});
	}

	function toggleBulb() {
		const map = {...light};
                map.Status = !map.Status;
                setLight(map);
		pushToggle(map.Status);
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

	function dropDown(e) {
		console.log("Dropdown", e);
	}

	function renderPatterns() {
	

		console.log("patterns", patterns);

		return (
			<div className="Popup Pattern">
				{ patterns.map(x=>{
					return <div onClick={()=>{ selectPattern(x.id)}}>{x.alias} - {x.creator}</div>
				}) }
			</div>
		);
	}

	function renderSettings() {
		return (
			<div className="Popup Auth">
				<LightSettings light={light} hub={hub} callback={settingsCallback}/>
			</div>
		);
	}

	function togglePatternMode() {
		setPatternMode(!patternMode);
	}

	return (
	<div className="Picker">
		{ settings ? renderSettings() : <></>}

		<div className="ButtonText" onClick={()=>setSettings(!settings)}>{light.Alias === "" ? light.ID : light.Alias}</div>
		<div className="Button" style={{background: color.hex}} onClick={toggleBulb}>
			<img className="ButtonImage" src={light.Status ? bulb : offbulb }></img>
		</div>
		<div>
			<p>{selected}</p>
		</div>
		{ light.Status ?  
			<>
			<input type="checkbox" onClick={()=>togglePatternMode()}/>
			{!patternMode ?
				<>
					<br/>
					<HuePicker className="ButtonImage" color={color.hsl} onChangeComplete={(c) => onColorChange(false, c)}/>
					<br/>
                			<AlphaPicker className="ButtonImage" color={color.hsl} onChangeComplete={(c) => onColorChange(true, c)}/>
				</>
			:
			renderPatterns()}
			</>

		: 
			<></>
		}

	</div>
  	);
}


export default Picker2;
