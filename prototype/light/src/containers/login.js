import React, { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { useFormFields, displayError } from "../lib/hooks";
import axios from "axios";

function Login({callback}) {

	const URL = process.env.REACT_APP_API_URL + "/user";

	const [field, handleFieldChange] = useFormFields({
		user: "",
    		password: "",
  	});

	const [loading, setLoading] = useState(false);

	function validateConfirmationForm() {
		return (
      			field.user.length > 0 &&
      			field.password.length > 0
    		);
	}

	function handleSubmit(event) {

		event.preventDefault();

		setLoading(true);

		axios.post(URL, JSON.stringify(field) ,{
    			withCredentials: true 
  		}).then((resp) => {
			callback()
		}).catch((e) => {
			displayError(e);
			setLoading(false);
		});
	}

	function renderForm() {
		return (
			<Form onSubmit={handleSubmit}>
				<Form.Group controlId="user" size="lg">
					<Form.Label>Username</Form.Label>
					<Form.Control
						autoFocus
						type="text"
						value={field.user}
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
				<br/>
				<Button
					block
					size="lg"
					type="submit"
					variant="success"
				        disabled={!validateConfirmationForm() || loading}

				>
					{ !loading ? "Login" : "loading..."}
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

export default Login;
