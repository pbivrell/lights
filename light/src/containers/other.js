import '../App.css';
import '../index.css';
import 'bootstrap/dist/css/bootstrap.css';
import React, { useState, useEffect } from "react";
import LightListItem from "../components/LightListItem";
import Picker from "../components/Picker";
import Styler from "../components/Styler";
import Offcanvas from "react-bootstrap/Offcanvas";
import Button from "react-bootstrap/Button";
import ListGroup from "react-bootstrap/ListGroup";
import axios from "axios";
import FileBrowser from 'react-keyed-file-browser'
import '../../node_modules/react-keyed-file-browser/dist/react-keyed-file-browser.css';


function Other() {

	async function getData() {


		/*return {
			lights: [
				{
					id: 0,
					name: "lower tree",
					ip: "192.168.86.146",
					pixels: 100,
				},
				{
					id: 1,
					name: "upper tree",
					ip: "192.168.86.145",
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
		}*/
	}

	const [data, setData] = useState({});

	function getData() {
		console.log("getting data");

		axios.get("http://homeserver/lights/config.json").then((resp) => {
			let blob = resp.data;
			setData(blob);
		});
	}

	const [files, setFiles] = useState([]);
	function getFiles() {
		axios.get("http://homeserver/lights/patterns/files.json").then((resp) => {
			let blob = resp.data;
			setFiles(blob);
		});

	}

	const [selectedFile, setSelectedFile] = useState();
	const [selectedLight, setSelectedLight] = useState();

	function clickFile(file) {

		var bodyFormData = new FormData();

		if (!selectedLight) {
			return;
		}

		axios.get(`http://homeserver/lights/patterns/${file.key}`,  { responseType: 'arraybuffer' }
).then((resp) => {

			let bytes = new Uint8Array(resp.data);

			let blob = new Blob([resp.data], {type: "application/octet-stream"});

			bodyFormData.append("file", blob, file.key);

			console.log(selectedLight);
			selectedLight.lights.forEach((light) => {
				axios({
					method: "post",
					url: `http://${data.lights[light].ip}/upload`,
					data: bodyFormData,
					headers: { "Content-Type": "multipart/form-data" },
				})
					.then(function (response) {
						//handle success
						console.log(response);
					})
					.catch(function (response) {
						//handle error
						console.log(response);
					});

				console.log(data.lights[light]);
			});
		});

	}

	function clickLight(light) {
		setSelectedLight(light);
		console.log(light);
	}

	useEffect(() => {
		getData()
		getFiles()
	}, []);


	const [show, setShow] = useState(false);

	const handleClose = () => setShow(false);
	const handleShow = () => setShow(true);

	return (
		<>

		<h3>Selected Light: <span>{selectedLight? selectedLight.name : "" }</span></h3>

		<Offcanvas show={show} onHide={handleClose} scroll={true} backdrop={false} placement={'bottom'}>
			<Offcanvas.Header closeButton>
				<Offcanvas.Title>Lights</Offcanvas.Title>
			</Offcanvas.Header>
			<Offcanvas.Body>
				<ListGroup>
				{
					data.groups ? (data.groups.map((v, i) => {
						return (
							<ListGroup.Item action key={i} onClick={ () => clickLight(v) } >
								<LightListItem name={v.name} all={data.lights} lights={v.lights}/>
							</ListGroup.Item>
						);
					})) : <></>
				}
				</ListGroup>
			</Offcanvas.Body>
		</Offcanvas>
		<FileBrowser files={files} onSelectFile={clickFile} />
		<Picker lightData={data} lights={selectedLight}/>
		<Styler/>
		<Button variant="primary" onClick={handleShow}>Lights </Button>
		</>
	);
}


export default Other;
