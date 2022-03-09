import React, { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { useFormFields, displayError } from "../lib/hooks";
import axios from "axios";

function Signup({callback}) {

	const URL = process.env.REACT_APP_API_URL + "/user";

	const [loading, setLoading] = useState(false);

	const [field, handleFieldChange] = useFormFields({
    		email: "",
		id: "",
    		password: "",
		confirmPassword: "",
    		confirmationCode: "",
  	});

	function validateConfirmationForm() {
		return (
      			field.email.length > 0 &&
      			field.password.length > 0 &&
      			field.password === field.confirmPassword
    		);
	}

	function handleSubmit(event) {

		event.preventDefault();

		setLoading(true);

		axios.put(URL, JSON.stringify(field)).then((resp) => {
			axios.post(URL, JSON.stringify({
				"user": field.id,
				"password": field.password,
			}) ,{
    				withCredentials: true 
  			}).then((resp) => {
				callback()
			}).catch((e) => {
				displayError(e);
				setLoading(false);
			});
		}).catch((e) => {
			setLoading(false);
			displayError(e);
		});
	}

	function renderForm() {
		return (
			<Form onSubmit={handleSubmit}>
				<Form.Group controlId="id" size="lg">
					<Form.Label>Username</Form.Label>
					<Form.Control
						autoFocus
						type="text"
						value={field.id}
						onChange={handleFieldChange}
					/>
				</Form.Group>
				<Form.Group controlId="email" size="lg">
					<Form.Label>Email</Form.Label>
					<Form.Control
						autoFocus
						type="email"
						value={field.email}
						onChange={handleFieldChange}
					/>
				</Form.Group>
				<Form.Group controlId="password" size="lg">
					<Form.Label>Password</Form.Label>
					<Form.Control
						type="password"
						value={field.password}
						onChange={handleFieldChange}
					/>
				</Form.Group>
				<Form.Group controlId="confirmPassword" size="lg">
					<Form.Label>Confirm Password</Form.Label>
					<Form.Control
						type="password"
						onChange={handleFieldChange}
						value={field.confirmPassword}
					/>
				</Form.Group>
				<br/>
				<Button
					block
					size="lg"
					type="submit"
					variant="success"
				        disabled={!validateConfirmationForm() || loading}

				>
					{ !loading ? "Signup" : "loading..." }
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

export default Signup;
