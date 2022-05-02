import React, { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { useFormFields, displayError } from "../lib/hooks";
import axios from "axios";
import './lightsettings.css';

function LightSettings({light, hub, callback}) {

	const URL = `${process.env.REACT_APP_API_URL}/user/hub/${hub}/light/${light.ID}`

	const [field, handleFieldChange] = useFormFields({
		Alias: light.Alias,
		Count: parseInt(light.Count),
  	});

	


	const [loading, setLoading] = useState(false);

	function validateConfirmationForm() {
		return (
      			field.Count > 0
    		);
	}

	function handleSubmit(event) {

		event.preventDefault();

		setLoading(true);

		field.Count = parseInt(field.Count);
		console.log(field);
		axios.post(URL, JSON.stringify(field) ,{
    			withCredentials: true 
  		}).then((resp) => {
			setLoading(false);
			callback();
		}).catch((e) => {
			displayError(e);
			setLoading(false);
		});
	}

	function renderForm() {
		return (
			<Form onSubmit={handleSubmit} className="SettingForm">
				<Form.Label>ID: {light.ID}</Form.Label><br/>
			        <Form.Group controlId="Alias" size="lg">
					<Form.Label>Alias:</Form.Label>
					<Form.Control
						autoFocus
						className="SettingsFormInput"
						type="text"
						value={field.Alias}
						onChange={handleFieldChange}
					/>
				</Form.Group>
			        <Form.Group controlId="Count" size="lg">
					<Form.Label>Bulbs:</Form.Label>
					<Form.Control
						autoFocus
						className="SettingsFormInput"
						type="number"
						value={parseInt(field.Count)}
						onChange={handleFieldChange}
					/>
				</Form.Group>
				<Form.Label>{new Date(light.Updated).toLocaleDateString() +" " +  new Date(light.Updated).toLocaleTimeString()}</Form.Label>
				<br/>
				<Button
					block
					size="sm"
					type="submit"
					variant="success"
				        disabled={!validateConfirmationForm() || loading}

				>
					{ !loading ? "Update" : "loading..."}
				</Button>
			</Form>
		);
	}

	return (
		<div className="Signup">
			{renderForm()}
		</div>
	);
}

export default LightSettings;
