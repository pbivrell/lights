(this.webpackJsonplight=this.webpackJsonplight||[]).push([[0],{202:function(e,t,n){"use strict";n.r(t);var c=n(0),a=n.n(c),i=n(23),r=n.n(i),s=(n(71),n(10)),o=(n(38),n(21)),l=n(29),j=n.n(l),h=n(36),u=n(42),b=n(20),d=n.n(b),O=n(95),p=n(1),f="http://192.168.1.115",x=3e4;function g(e){var t=e.next,n=Object(c.useState)(!1),a=Object(s.a)(n,2),i=a[0],r=a[1],l=Object(c.useState)(!1),b=Object(s.a)(l,2),g=b[0],m=b[1],v=Object(c.useState)([]),y=Object(s.a)(v,2),w=y[0],k=y[1],S=Object(c.useState)(""),C=Object(s.a)(S,2),F=C[0],A=C[1],N=Object(c.useState)(""),T=Object(s.a)(N,2),B=T[0],E=T[1];function D(){return L.apply(this,arguments)}function L(){return(L=Object(h.a)(j.a.mark((function e(){return j.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:r(!0),d.a.get(f+"/networks",{timeout:x}).then((function(e){k(e.data)})).catch((function(e){m(!0)})),r(!1);case 3:case"end":return e.stop()}}),e)})))).apply(this,arguments)}function P(){return(P=Object(h.a)(j.a.mark((function e(){return j.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:r(!0),d.a.get(f+"/login?ssid=".concat(F,"&password=").concat(B),{timeout:x}).then((function(e){t()})).catch((function(e){m(!0)})),r(!1);case 3:case"end":return e.stop()}}),e)})))).apply(this,arguments)}return Object(c.useEffect)((function(){D()}),[]),Object(p.jsxs)("div",{children:[i?Object(p.jsxs)(p.Fragment,{children:[Object(p.jsx)("p",{children:" Searching for lights "}),Object(p.jsx)(u.a,{animation:"border"}),Object(p.jsx)("br",{})]}):g?Object(p.jsxs)(p.Fragment,{children:[Object(p.jsx)("p",{children:"Failed to setup new lights"}),Object(p.jsx)(o.a,{onClick:function(){return m(!1),void D()},children:"Retry"})]}):Object(p.jsxs)(p.Fragment,{children:[Object(p.jsx)("p",{children:"Found new lights"}),Object(p.jsx)("p",{children:"Select Wifi you would like to connect them to"}),Object(p.jsx)(O.a,{onChange:function(e){return A(e.value)},searchable:!0,options:w.map((function(e){return{value:e,label:e}}))}),Object(p.jsx)("label",{children:"Password: "})," ",Object(p.jsx)("input",{type:"password",name:"name",onChange:function(e){return E(e.target.value)}}),Object(p.jsx)("br",{}),Object(p.jsx)(o.a,{onClick:function(){return P.apply(this,arguments)},children:"Submit"})]}),Object(p.jsx)(o.a,{onClick:function(){t(-1)},children:"Back"})]})}var m="http://192.168.1.115";function v(e){var t=e.next,n=Object(c.useState)(!1),a=Object(s.a)(n,2),i=a[0],r=a[1],l=Object(c.useState)(!1),b=Object(s.a)(l,2),O=b[0],f=b[1],x=Object(c.useState)(""),g=Object(s.a)(x,2),v=(g[0],g[1]);function y(){return w.apply(this,arguments)}function w(){return(w=Object(h.a)(j.a.mark((function e(){return j.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:r(!0),d.a.get(m+"/ip",{timeout:3e3}).then((function(e){v(e.data)})).catch((function(e){f(!0)})),r(!1);case 3:case"end":return e.stop()}}),e)})))).apply(this,arguments)}return Object(c.useEffect)((function(){y()}),[]),Object(p.jsxs)("div",{children:[i?Object(p.jsxs)(p.Fragment,{children:[Object(p.jsx)("p",{children:" Getting device IP"}),Object(p.jsx)(u.a,{animation:"border"}),Object(p.jsx)("br",{})]}):O?Object(p.jsxs)(p.Fragment,{children:[Object(p.jsx)("p",{children:"Failed to connect with device"}),Object(p.jsx)(o.a,{onClick:function(){y()},children:"Retry"})]}):Object(p.jsxs)(p.Fragment,{children:[Object(p.jsx)("p",{children:"Congradulations your lights are ready to go!"}),Object(p.jsx)("p",{children:"Remeber to reconnect your device to the same network as the lights"}),Object(p.jsx)(o.a,{onClick:function(){return console.log("TODO")},children:"Light"})]}),Object(p.jsx)(o.a,{onClick:function(){t(-1)},children:"Back"})]})}var y=function(){var e=Object(c.useState)(0),t=Object(s.a)(e,2),n=t[0],a=t[1];function i(){var e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:1,t=n+e;t>3&&(t=0),a(t)}return Object(p.jsx)("div",{className:"App",children:Object(p.jsx)("header",{className:"App-header",children:Object(p.jsx)("div",{children:0===n?Object(p.jsxs)(p.Fragment,{children:[Object(p.jsxs)("p",{children:["Plug in Inf-Lights. They should blink ",Object(p.jsx)("span",{class:"red",children:"red"})]}),Object(p.jsx)(o.a,{onClick:function(){i()},children:" Next"})]}):1===n?Object(p.jsxs)(p.Fragment,{children:[Object(p.jsxs)("p",{children:["Go to wifi settings and connect to ",Object(p.jsx)("code",{children:"P-lights-0"})]}),Object(p.jsx)(o.a,{onClick:function(){i()},children:" Next"}),Object(p.jsx)(o.a,{onClick:function(){return i(-1)},children:" Back"})]}):2===n?Object(p.jsx)(g,{next:i}):3===n?Object(p.jsx)(v,{next:i}):Object(p.jsx)(p.Fragment,{})})})})},w=n(91);var k=function(){var e=Object(c.useState)({r:200,g:150,b:35}),t=Object(s.a)(e,2),n=t[0],a=t[1],i=Object(c.useState)("192.168.86.126"),r=Object(s.a)(i,2),o=r[0],l=(r[1],Date.now());return Object(p.jsx)("div",{className:"App",children:Object(p.jsxs)("header",{className:"App-header",children:[Object(p.jsx)(w.a,{color:n,onChange:function(e){var t=Date.now();console.log(e),t-l<100?console.log("ignoring time"):(l=t,a(e),d.a.get("http://".concat(o,"/color?r=").concat(e.r,"&g=").concat(e.g,"&b=").concat(e.b)).catch((function(e){console.log("TODO")})))}}),Object(p.jsxs)("p",{children:[n.r,", ",n.g,", ",n.b]})]})})},S=n(40);function C(e){var t=e.delay,n=e.color,a=e.x,i=e.y,r=e.radius,o=Object(c.useState)(!1),l=Object(s.a)(o,2),j=l[0],h=l[1];return Object(c.useEffect)((function(){setTimeout((function(){h(!0)}),t)}),[t]),j?Object(p.jsx)(S.a,{x:a,y:i,radius:r,fill:n}):Object(p.jsx)(S.a,{x:a,y:i,radius:r,fill:"black"})}var F=function(){function e(t){console.log(t);var n=new Uint8Array(t);console.log(n);for(var c=[],a=0,r=0;r<n.length;r+=8)if(3===n[r]&&0===r){var s=n[r+2]<<8|n[r+1];console.log("Count",s);for(var o=0;o<s;o++)c.push({color:"rgb(0,0,0)"})}else if(1===n[r]){var l=n[r+1],j=n[r+2],h=n[r+3],u=n[r+5]<<8|n[r+4];console.log("color",l,j,h,u),c[u]={color:"rgb(".concat(l,",").concat(j,",").concat(h,")")}}else 2===n[r]&&(a+=n[r+4]<<24|n[r+3]<<16|n[r+2]<<8|n[r+1],console.log("delay",a),setTimeout((function(e){console.log("Changed lights"),i(e)}),a,[].concat(c)));setTimeout((function(){e(t)}),a)}var t=Object(c.useState)([]),n=Object(s.a)(t,2),a=n[0],i=n[1],r=Math.floor((window.innerWidth-200)/170);return Object(p.jsxs)("div",{className:"App",children:[Object(p.jsx)("input",{type:"file",name:"file",onChange:function(t){var n=new FileReader;n.onload=function(){e(n.result)},t.target.files.length>0&&n.readAsArrayBuffer(t.target.files[0])}}),Object(p.jsx)("header",{className:"App-header",children:Object(p.jsx)(S.c,{width:window.innerWidth,height:window.innerHeight,children:Object(p.jsx)(S.b,{children:a.map((function(e,t){return console.log(e,t),Object(p.jsx)(C,{x:120+t%r*120,y:120+120*Math.floor(t/r),color:e.color,radius:50})}))})})})]})};var A=function(){var e=Object(c.useState)("192.168.86.152"),t=Object(s.a)(e,2),n=t[0];return t[1],Object(p.jsx)("div",{className:"App",children:Object(p.jsx)("header",{className:"App-header",children:Object(p.jsxs)("form",{method:"POST",enctype:"multipart/form-data",action:"http://"+n+"/upload",children:[Object(p.jsx)("input",{type:"file",id:"myFile",name:"filename"}),Object(p.jsx)("input",{type:"submit"})]})})})};function N(e){var t=e.name,n=e.all,a=e.lights,i=Object(c.useState)(!1),r=Object(s.a)(i,2),o=r[0],l=r[1];Object(c.useEffect)((function(){1===a.length&&d.a.get("http://".concat(n[a[0]].ip,"/status"),{timeout:1e3}).then((function(e){console.log(e.data.status),e.data.status?b("active"):b("inactive")})).catch((function(){b("offline")}))}),[o]);var j=Object(c.useState)(1===a.length?"loading":""),h=Object(s.a)(j,2),u=h[0],b=h[1];return Object(p.jsxs)("div",{style:"offline"===u?{color:"grey"}:{},children:[Object(p.jsx)("input",{type:"checkbox",onClick:function(e){return a.forEach((function(e){d.a.get("http://".concat(n[e].ip,"/toggle"),{timeout:1e3}).catch((function(e){console.log("failed to toggle")}))})),void setTimeout((function(){l(!o)}),500)},defaultChecked:"active"===u,style:{float:"right",width:"100px"}}),Object(p.jsx)("span",{children:"Name: "}),Object(p.jsx)("span",{children:t}),Object(p.jsx)("br",{}),Object(p.jsx)("span",{children:"Devices(s): "}),Object(p.jsx)("span",{children:1===a.length?n[a[0]].ip:a.map((function(e){return n[e].name})).join(",")})]},t)}var T=function(){var e={lights:[{id:0,name:"lower tree",ip:"192.168.86.136",pixels:100},{id:1,name:"upper tree",ip:"192.168.86.137",pixels:100},{id:2,name:"test",ip:"192.168.86.140",pixels:50}],groups:[{name:"lower tree",lights:[0]},{name:"upper tree",lights:[1]},{name:"test",lights:[2]},{name:"tree",lights:[0,1]}]};return Object(p.jsx)("div",{className:"App",children:Object(p.jsx)("header",{className:"App-header",children:Object(p.jsx)("ul",{children:e.groups.map((function(t,n){return Object(p.jsx)("li",{children:Object(p.jsx)(N,{name:t.name,all:e.lights,lights:t.lights})})}))})})})},B=(n(139),n(44)),E=n(70),D=n(93),L=n.n(D);n(196);var P=function(){var e=Object(c.useState)({}),t=Object(s.a)(e,2),n=t[0],a=t[1],i=Object(c.useState)([]),r=Object(s.a)(i,2),l=r[0],j=r[1],h=Object(c.useState)(),u=Object(s.a)(h,2),b=(u[0],u[1],Object(c.useState)()),O=Object(s.a)(b,2),f=O[0],x=O[1];Object(c.useEffect)((function(){console.log("getting data"),d.a.get("http://homeserver/lights/config.json").then((function(e){var t=e.data;a(t)})),d.a.get("http://homeserver/lights/patterns/files.json").then((function(e){var t=e.data;j(t)}))}),[]);var g=Object(c.useState)(!1),m=Object(s.a)(g,2),v=m[0],y=m[1];return Object(p.jsxs)(p.Fragment,{children:[Object(p.jsx)(o.a,{variant:"primary",onClick:function(){return y(!0)},children:"Launch"}),Object(p.jsxs)("h3",{children:["Selected Light: ",Object(p.jsx)("span",{children:f?f.name:""})]}),Object(p.jsxs)(B.a,{show:v,onHide:function(){return y(!1)},scroll:!0,backdrop:!1,children:[Object(p.jsx)(B.a.Header,{closeButton:!0,children:Object(p.jsx)(B.a.Title,{children:"Lights"})}),Object(p.jsx)(B.a.Body,{children:Object(p.jsx)(E.a,{children:n.groups?n.groups.map((function(e,t){return Object(p.jsx)(E.a.Item,{action:!0,onClick:function(){return x(t=e),void console.log(t);var t},children:Object(p.jsx)(N,{name:e.name,all:n.lights,lights:e.lights})},t)})):Object(p.jsx)(p.Fragment,{})})})]}),Object(p.jsx)(L.a,{files:l,onSelectFile:function(e){var t=new FormData;f&&d.a.get("http://homeserver/lights/patterns/".concat(e.key),{responseType:"arraybuffer"}).then((function(c){new Uint8Array(c.data);var a=new Blob([c.data],{type:"application/octet-stream"});t.append("file",a,e.key),console.log(f),f.lights.forEach((function(e){d()({method:"post",url:"http://".concat(n.lights[e].ip,"/upload"),data:t,headers:{"Content-Type":"multipart/form-data"}}).then((function(e){console.log(e)})).catch((function(e){console.log(e)})),console.log(n.lights[e])}))}))}})]})},I=function(e){e&&e instanceof Function&&n.e(3).then(n.bind(null,207)).then((function(t){var n=t.getCLS,c=t.getFID,a=t.getFCP,i=t.getLCP,r=t.getTTFB;n(e),c(e),a(e),i(e),r(e)}))},R=n(66),H=n(15);r.a.render(Object(p.jsx)(a.a.StrictMode,{children:Object(p.jsx)(R.a,{children:Object(p.jsxs)(H.c,{children:[Object(p.jsx)(H.a,{exact:!0,path:"/dev",children:Object(p.jsx)(F,{})}),Object(p.jsx)(H.a,{exact:!0,path:"/run",children:Object(p.jsx)(A,{})}),Object(p.jsx)(H.a,{exact:!0,path:"/setup",children:Object(p.jsx)(y,{})}),Object(p.jsx)(H.a,{exact:!0,path:"/lights",children:Object(p.jsx)(k,{})}),Object(p.jsx)(H.a,{exact:!0,path:"/control",children:Object(p.jsx)(T,{})}),Object(p.jsx)(H.a,{exact:!0,path:"/other",children:Object(p.jsx)(P,{})})]})})}),document.getElementById("root")),I()},38:function(e,t,n){},71:function(e,t,n){}},[[202,1,2]]]);
//# sourceMappingURL=main.cf4e16ea.chunk.js.map