import { useState } from "react";

export function useFormFields(initialState) {
  const [fields, setValues] = useState(initialState);

  return [
    fields,
    function(event) {
      setValues({
        ...fields,
        [event.target.id]: event.target.value
      });
    }
  ];
}

export function displayError(e) {
	if (e.response) {

		console.log("URL", e.request.responseURL);	
		if (e.response.status === 401 && e.request.responseURL !==  "https://lights.paulbivrell.com/user") {
			console.log("reload please");

		}
	}
	console.log(e);
}
