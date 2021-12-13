import React, { useState, useEffect } from "react";
import axios from "axios";

export default function LightDetect({name, all, lights}) {

	const [clicked, setClicked] = useState(false);
	
	function onToggle(e) {
		lights.forEach((v) => {
			axios.get(`http://${all[v].ip}/toggle`, {timeout: 1000}).catch((resp) => {
				console.log("failed to toggle");
			});
		});

		// After 1/2 a second check if the lights are on
		setTimeout(() => {
			setClicked(!clicked);
		}, 500);
	}


	useEffect(() => {
		if (lights.length === 1) {
			axios.get(`http://${all[lights[0]].ip}/status`, {timeout: 1000})
			.then((resp) => {
				console.log(resp.data.status);
				if (resp.data.status) {
					setStatus("active");
				}else {
					setStatus("inactive");
				}
			}).catch(() => {
				setStatus("offline");
			});
		}
	}, [clicked])


	const [Status, setStatus] = useState(lights.length === 1 ? "loading" : "");

	return (
    	<div key={name} style={ Status === "offline" ? {color: "grey"} : {}}>
		<input
          		type='checkbox'
			onClick={(e) => onToggle(e)}
			defaultChecked={Status === "active"}
			style={{float: "right", width: "100px"}}
        	/>
		<span>Name: </span><span>{name}</span><br/>
		<span>Devices(s): </span><span>{ lights.length === 1 ? all[lights[0]].ip : lights.map((x) => all[x].name ).join(",") }</span>
    	</div>
  	);
}
