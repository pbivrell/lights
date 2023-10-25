import '../App.css';
import '../index.css';
import './home.css';
import '../components/Picker2.css'
import Login from "./login.js";
import Signup from "./signup.js";
import {displayError} from "../lib/hooks";
import React, { useState, useEffect } from "react";
import axios from "axios";
import {ToggleButton, ToggleButtonGroup, DropdownButton, Dropdown, ListGroup} from "react-bootstrap";
import Picker2 from "../components/Picker2.js"
import AddBulb from "../res/AddBulb_crop.png";
import Setup from '../containers/setup.js';

function Home() {

	const [setup, setSetup] = useState(false);

	const [hubs, setHubs] = useState([]);

	const [userData, setUserData] = useState({});

	const [authed, setAuthed] = useState(true);

	const [register, setRegister] = useState(false);

	function addBulbToggle() {
		setSetup(!setup);
	}

	function dropDown(e) {

		if (e === "logout") {
			logout()
		}else if(e === "settings") {
			settings()
		}
	}

	function logout() {
		axios.delete(`${process.env.REACT_APP_API_URL}/user`, {
			withCredentials: true,
		}).then((resp) => {
			setAuthed(false);
			setUserData({});
		}).catch((e) => {
			displayError(e);
		});
		console.log("logout");
	}

	function settings() {
		console.log("setting");
	}

	function sampleUserData() {
		return {
			"id": "ruckens",
			"lights": [
				{
					"alias": "main",
					"id": "OriginalLight1"
				},
				{
					"alias": "main",
					"id": "OriginalLight1"
				}
			],
		}
	}

	function getUserData() {
		
		console.log("Loading", authed);
		axios.get(`${process.env.REACT_APP_API_URL}/user`, {
			withCredentials: true,
		}).then((resp) => {
			console.log("Test",resp);
			loadUserData(resp.data);
		}).catch((e) => {
			setAuthed(false);
			displayError(e);
		});
	}

	function loadUserData(data) {
		setUserData(data);
		
		if (data.hubs) {
			data.hubs.forEach((v) => {
				lookupHub(v.id);
			});
		}

	}

	function lookupHub(name) {
		axios.get(`${process.env.REACT_APP_API_URL}/user/hub/${name}`, {
			withCredentials: true,
		}).then((resp) => {
			loadHubData(name, resp.data);
		}).catch((e) => {
			displayError(e);
		});
	}

	function loadHubData(hub, data) {

		const map = {...hubs};
		map[hub] = data;
		setHubs(map);
	}

	function renderRegister() {
		return (
			<div className="Auth">
				<Setup/>
			</div>
		);
	}

	function renderAuth() {
		return (
			<div className="Popup Centered">
				<ToggleButtonGroup type="checkbox" value={register} onChange={(v) => setRegister(!register) }>
      					<ToggleButton id="tbg-btn-1" value={false}>
      						Login	
					</ToggleButton>
      					<ToggleButton id="tbg-btn-3" value={true}>
      						Signup
					</ToggleButton>
    				</ToggleButtonGroup>
				{ !register ?
					<>
					<h3>Login</h3>
					<Login callback={ () => {
						setAuthed(true);
						console.log("Logged in");
						window.location.reload();
					}}/>
					</>
				:
					<>
					<h3>Signup</h3>
					<Signup callback={ () => setAuthed(true) }/>
					</>
					
				}
			</div>
		)

	}

	function renderPage() {
		return (
			<div className="Page">
				<div className="User">
			      		<DropdownButton
						
        					drop={'up'}
        					variant="secondary"
						title={userData.id ? userData.id : ""}
						size="lg"
						onSelect={dropDown}
      					>
        					<Dropdown.Item eventKey="settings">Settings</Dropdown.Item>
        					<Dropdown.Divider />
        					<Dropdown.Item eventKey="logout">Logout</Dropdown.Item>
      					</DropdownButton>
				</div>
				<div className="ButtonPanel">
					{ hubs !== null ? Object.entries(hubs).map(([key, value])=>{
						console.log("This", hubs);
						if (value.lights.length < 1 ) {
							return <><p>None of your lights are on! </p></>
						}else {
							return value.lights.map((value)=>{
								return <Picker2 className="ButtonItem" hub={key} name={value}/>
							});
						}
					}) : 
						<>
							<p>You don't appear to have a registed hub <a href="#">register now</a>. If this is a mistake get <a href="#">more info</a></p>
						</>
					}
				</div>
			</div>
		);
	}

	useEffect(() => {
		getUserData()
		setTimeout(()=>{
			setAuthed(false);
		}, 5 * 60 * 1000)
		//console.log(sampleUserData().lights);
		//setUserData(sampleUserData());
	}, []);


	return (
		<>
			{ authed ? 
				renderPage()
			:
				renderAuth()
			}

		</>
	);
}


export default Home;
