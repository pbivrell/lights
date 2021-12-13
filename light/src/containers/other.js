import '../App.css';
import '../index.css';
import 'bootstrap/dist/css/bootstrap.css';
import React, { useState, useEffect } from "react";
import LightListItem from "../components/LightListItem";
import Offcanvas from "react-bootstrap/Offcanvas";
import Button from "react-bootstrap/Button";
import ListGroup from "react-bootstrap/ListGroup";
import axios from "axios";

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

	useEffect(() => {
		axios.get("http://homeserver/lights/config.json").then((resp) => {
			let blob = resp.data;
			setData(blob);
		});
	}, []);


	const [show, setShow] = useState(false);

	const handleClose = () => setShow(false);
  	const handleShow = () => setShow(true);

 	 return (
		<>
      		<Button variant="primary" onClick={handleShow}>
        		Launch
      		</Button>

      		<Offcanvas show={show} onHide={handleClose} scroll={true} backdrop={false}>
        		<Offcanvas.Header closeButton>
          			<Offcanvas.Title>Lights</Offcanvas.Title>
        		</Offcanvas.Header>
        		<Offcanvas.Body>
	  				<ListGroup>
			 				{
								data.groups ? 
								(data.groups.map((v, i) => {
									return (
										<ListGroup.Item action key={i}>
                                			<LightListItem name={v.name} all={data.lights} lights={v.lights}/>
			 							</ListGroup.Item>
									);
								}))
								: <></>
							}
					</ListGroup>
        		</Offcanvas.Body>
      		</Offcanvas>
    	</>
  	);
}


export default Other;
