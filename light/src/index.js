import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import Setup from './containers/setup.js';
import Light from './containers/light.js';
import Programming from './containers/programming.js';
import Run from './containers/run.js';
import Control from './containers/control.js';
import Other from './containers/other.js';
import reportWebVitals from './reportWebVitals';
import { BrowserRouter as Router } from 'react-router-dom';
import { Route, Switch } from "react-router-dom";

ReactDOM.render(
  <React.StrictMode>
    <Router>
      <Switch>
      	<Route exact path="/dev">
          <Programming />
      	</Route>
      	<Route exact path="/run">
          <Run />
      	</Route>
      	<Route exact path="/setup">
          <Setup />
      	</Route>
      	<Route exact path="/lights">
          <Light/>
      	</Route>
      	<Route exact path="/control">
          <Control />
      	</Route>
      	<Route exact path="/">
          <Other/>
      	</Route>
      </Switch>
    </Router>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
