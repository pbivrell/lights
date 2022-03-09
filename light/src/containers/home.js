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

	const [lightMap, setLightMap] = useState({});

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
		
		axios.get(`${process.env.REACT_APP_API_URL}/user`, {
			withCredentials: true,
		}).then((resp) => {
			loadUserData(resp.data);
		}).catch((e) => {
			setAuthed(false);
			displayError(e);
		});
	}

	function loadUserData(data) {
		setUserData(data);
		
		data.lights.forEach((v) => {
			lookupLight(v.id);
		});

	}

	function lookupLight(name) {
		
		axios.get(`${process.env.REACT_APP_API_URL}/light?light=${name}`, {
			withCredentials: true,
		}).then((resp) => {
			let map = {...lightMap};
			let lightData = resp.data;

			map[name] = lightData;
			setLightMap(map)
			console.log(lightMap);
		}).catch((e) => {
			displayError(e);
		});
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
			<div className="Auth">
				<ToggleButtonGroup type="checkbox" value={register} onChange={(v) => setRegister(!register) }>
      					<ToggleButton id="tbg-btn-1" value={true}>
      						Login	
					</ToggleButton>
      					<ToggleButton id="tbg-btn-3" value={false}>
      						Signup
					</ToggleButton>
    				</ToggleButtonGroup>
				{ register ?
					<>
					<h3>Login</h3>
					<Login callback={ () => {
						setAuthed(true) 
						getUserData()
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
						title={userData.id}
						size="lg"
						onSelect={dropDown}
      					>
        					<Dropdown.Item eventKey="settings">Settings</Dropdown.Item>
        					<Dropdown.Divider />
        					<Dropdown.Item eventKey="logout">Logout</Dropdown.Item>
      					</DropdownButton>
				</div>
				{ setup ? <Setup cancel={addBulbToggle}/> : <></> }
				<div className="ButtonPanel">
					{ userData.lights ? userData.lights.map((v) => {
						return <Picker2 light={lightMap[v.id]} className="ButtonItem" initalState={1} name={v}/>
					}) : <></>
					}
					<div className="Picker" onClick={addBulbToggle}>
                        			<img className="ButtonImage" src={AddBulb}></img>
					</div>
				</div>
			</div>
		);
	}

	useEffect(() => {
		getUserData()
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
