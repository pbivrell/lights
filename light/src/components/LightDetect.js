import React, { useState, useEffect } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Spinner from "react-bootstrap/Spinner";
import axios from "axios";
import Select from 'react-select'


const baseURL = process.env.REACT_APP_URL;
const timeout = 30 * 1000;

export default function LightDetect({next}) {
	const [loading, setLoading] = useState(false);
	const [failed, setFailed] = useState(false);
	const [networks, setNetworks] = useState([]);
	const [ssid, setSSID] = useState("");
	const [password, setPassword] = useState("");


	useEffect(() => {
  		onLoad();
	}, []);

	async function onLoad() {
		setLoading(true);
		axios.get(baseURL + "/networks", { timeout: timeout}).then((response) => {
			setNetworks(response.data);
    		}).catch((resp) => {
			setFailed(true);
		});
		setLoading(false);
	}

	async function onSubmit() {
		setLoading(true);

		axios.get(baseURL + `/login?ssid=${ssid}&password=${password}`, { timeout: timeout}).then((response) => {
			next();
    		}).catch((resp) => {
			setFailed(true);
		});
		setLoading(false);

	}

	function onRetry() {
		setFailed(false);
		onLoad();
	}

	return (
    	<div>
			{ loading ? (
				<>
					<p> Searching for lights </p>
					<Spinner animation="border">
					</Spinner>
					<br/>
				</>
			) : failed ? (
				<>
					<p>Failed to setup new lights</p>
					<Button onClick={() => onRetry()}>Retry</Button> 
				</>
			) : (
				<>
					<p>Found new lights</p>
					<p>Select Wifi you would like to connect them to</p>
					<Select onChange={(selected) => setSSID(selected.value)} searchable options={networks.map(function(d){
         					return ({ "value": d, "label": d })
       					})}/>
					<label>Password: </label> <input type="password" name="name" onChange={(e) => setPassword(e.target.value)}/><br/>
					<Button onClick={onSubmit}>Submit</Button>	
			    	</> 
			) 
			}
			<Button onClick={() => {next(-1)}}>Back</Button>
    	</div>
  	);
}
