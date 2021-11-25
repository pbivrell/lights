import React, { useState, useEffect } from "react";
import {Circle}  from 'react-konva';

export default function Bulb({delay, color, x, y, radius}) {


	const [isShown, setIsShown] = useState(false);

  	useEffect(() => {
    		setTimeout(() => {
      			setIsShown(true);
    		}, delay);
  	}, [delay]);

	return isShown ? 
		<Circle x={x} y={y} radius={radius} fill={color} /> 
		: 
		<Circle x={x} y={y} radius={radius} fill="black" />;
}
