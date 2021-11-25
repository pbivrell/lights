import React, { useState, useEffect } from "react";
import Button from "react-bootstrap/Button";
import axios from "axios";
import Spinner from "react-bootstrap/Spinner";

const baseURL = process.env.REACT_APP_URL;
const timeout = 3 * 1000;

export default function LightDetect({next}) {
	const [loading, setLoading] = useState(false);
	const [failed, setFailed] = useState(false);
	const [ip, setIP] = useState("");


	useEffect(() => {
  		onLoad();
	}, []);

	async function onLoad() {
		setLoading(true);
		axios.get(baseURL + "/ip", { timeout: timeout}).then((response) => {
			setIP(response.data);
    		}).catch((resp) => {
			setFailed(true);
		});
		setLoading(false);
	}

	function onRetry() {
		onLoad();
	}

	return (
    	<div>
			{ loading ? (
				<>
					<p> Getting device IP</p>
					<Spinner animation="border">
					</Spinner>
					<br/>
				</>
			) : failed ? (
				<>
					<p>Failed to connect with device</p>
					<Button onClick={() => onRetry()}>Retry</Button> 
				</>
			) : (
				<>
					<p>Congradulations your lights are ready to go!</p>
					<p>Remeber to reconnect your device to the same network as the lights</p>
					<Button onClick={() => console.log("TODO") }>Light</Button>
			    	</> 
			) 
			}
			<Button onClick={() => {next(-1)}}>Back</Button>
    	</div>
  	);
}
