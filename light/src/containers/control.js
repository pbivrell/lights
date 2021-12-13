import '../App.css';
import React, { useState, useEffect } from "react";
import LightListItem from "../components/LightListItem";

function Control() {

	const sampleData = {
		lights: [
			{
				id: 0,
				name: "lower tree",
				ip: "192.168.86.136",
				pixels: 100,
			},
			{
				id: 1,
				name: "upper tree",
				ip: "192.168.86.137",
				pixels: 100,
			},
			{
				id: 2,
				name: "test",
				ip: "192.168.86.140",
				pixels: 50,
			}
		],
		groups: [
			{
				name: "lower tree",
				lights: [0],
			},
			{
				name: "upper tree",
				lights: [1],
			},
			{
				name: "test",
				lights: [2],
			},	
			{
				name: "tree",
				lights: [ 0, 1 ],
			}
		]
	}

	return (
	<div className="App">
		<header className="App-header">
			<ul>
				{
					sampleData.groups.map((v,i) => {
						return (
							<li>
							<LightListItem name={v.name} all={sampleData.lights} lights={v.lights}/>
							</li>
						);
					})
				}
			</ul>
		</header>
	</div>
  	);
}


export default Control;
