import logo from './logo.svg';
import './App.css';
import Button from "react-bootstrap/Button";
import React, { useState } from "react";
import LightsDetect from "./components/LightDetect";
import LightsIP from "./components/LightIP";

function App() {
	const [prom, setPrompt] = useState(0);
	
	function updatePrompt(i = 1) {
		console.log("S",prom)
		let nextProm = prom + i;
		if (nextProm > 3) {
			nextProm = 0
		}
		console.log("E", nextProm);
   		setPrompt(nextProm);
	}


	return (
    	<div className="App">
      		<header className="App-header">
        		<div>
				{ prom === 0 ? (
					<>
						<p>Plug in Inf-Lights. They should blink <span class="red">red</span></p>
						<Button onClick={() => {updatePrompt()}}> Next</Button>
					</>
				)
				: prom === 1 ? (
					<>
						<p>Go to wifi settings and connect to <code>P-lights-0</code></p>
						<Button onClick={() => {updatePrompt()}}> Next</Button>
						<Button onClick={() => updatePrompt(-1)}> Back</Button>
					</>
				) : prom === 2 ? (
					<LightsDetect next={updatePrompt}/>
				) : prom === 3 ? (
					<LightsIP next={updatePrompt}/>
				) : (<></>) }
        		</div>
      		</header>
    	</div>
  	);
}


export default App;
