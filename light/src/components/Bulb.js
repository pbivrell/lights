import React, { useState, useEffect } from "react";
import {Circle}  from 'react-konva';

export default function Bulb({delay, color, x, y, radius, click}) {


	const [isShown, setIsShown] = useState(false);

  	useEffect(() => {
    		setTimeout(() => {
      			setIsShown(true);
    		}, delay);
  	}, [delay]);

	return isShown ? 
		<Circle x={x} y={y} radius={radius} fill={color} onClick={click}/> 
		: 
		<Circle x={x} y={y} radius={radius} fill="black" onClick={click}/>;
}
