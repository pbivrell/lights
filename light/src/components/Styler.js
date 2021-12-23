import '../App.css';
import React, { useState, useEffect } from "react";
import { Stage, Layer } from 'react-konva';
import Bulb from "../components/Bulb.js"

function Styler() {

	function processSequence(seq) {
		
		console.log(seq);

		var bytes = new Uint8Array(seq);

		console.log(bytes);

		let nextColors = [];
		
		let delay = 0;
			
		for (let i=0; i<bytes.length; i+=8) {
			if (bytes[i] === 0x3 && i === 0) {
				let count = (bytes[i+2] << 8) | bytes[i+1];
				console.log("Count", count);
				for (let j = 0; j < count; j++) {
					nextColors.push({color: `rgb(0,0,0)`})
				}
				
			} else if (bytes[i] === 0x1) {
				let r = bytes[i+1];
				let g = bytes[i+2];
				let b =  bytes[i+3];
				let index = (bytes[i+5] << 8) | bytes[i+4];
				console.log("color", r,g,b, index);
				nextColors[index] = ({color: `rgb(${r},${g},${b})`})
			}else if (bytes[i] === 0x2) {
				delay += (bytes[i+4] << 24) | (bytes[i+3] << 16) | (bytes[i+2] << 8) | bytes[i+1];
				console.log("delay", delay);
				setTimeout((lights) => {
					console.log("Changed lights");
					setLights(lights);
				}, delay,[...nextColors] )
			}
		}
        }

	const [lights, setLights] = useState([{color: `rgb(50,50,50)`},{color: `rgb(50,50,50)`}, {color: `rgb(50,50,50)`}, {color: `rgb(50,50,50)`}, {color: `rgb(50,50,50)`}, {color: `rgb(50,50,50)`}, {color: `rgb(50,50,50)`}, {color: `rgb(50,50,50)`}, {color: `rgb(50,50,50)`}, {color: `rgb(50,50,50)`} ]);
	const rows = Math.floor((window.innerWidth - 200) / 170);
	const radius = 50;
	const xOffset = 120;
	const yOffset = 120;

	return (
		<Stage width={window.innerWidth} height={200}>
      			<Layer>
				{
					lights.map( function(elem, idx) {
						console.log(elem, idx);
						return <Bulb x={xOffset + (xOffset * (idx % rows))} y={yOffset + (yOffset * Math.floor(idx / rows))} color={elem.color} radius={radius}/>
					})
				}
      			</Layer>
    		</Stage>
  	);
}


export default Styler;
