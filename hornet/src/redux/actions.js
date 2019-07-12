import {
  ADD_NODE_UPDATE,
  ADD_REGION,
  CONNECT,
  CONNECTED,
  CONNECTION_ERROR,
  DISCONNECTED,
  UPDATE_CPU_READ,
  UPDATE_DISK_READ,
  UPDATE_LAST_BLOCK_INFO,
  UPDATE_LOG_READ,
  UPDATE_MEM_READ,
  UPDATE_NET_READ,
  UPDATE_TIME_READ,
  UPDATE_WARN_LIST,
} from "./action-types";

export const addNodeUpdate = payload => ({
  type: ADD_NODE_UPDATE,
  payload,
});

export const updateCPURead = (value, timestamp) => ({
  type: UPDATE_CPU_READ,
  value,
  timestamp,
});

export const updateNetRead = (value, timestamp) => ({
  type: UPDATE_NET_READ,
  value,
  timestamp,
});

export const updateDiskRead = (value, timestamp) => ({
  type: UPDATE_DISK_READ,
  value,
  timestamp,
});

export const updateLogRead = (value, timestamp) => ({
  type: UPDATE_LOG_READ,
  value,
  timestamp,
});

export const updateMemoryRead = (value, timestamp) => ({
  type: UPDATE_MEM_READ,
  value,
  timestamp,
});

export const updateBlockTimeRead = (value, timestamp) => ({
  type: UPDATE_TIME_READ,
  value,
  timestamp,
});

export const updateWarningList = (value, timestamp) => ({
  type: UPDATE_WARN_LIST,
  value,
  timestamp,
});

export const addRegion = payload => ({
  type: ADD_REGION,
  payload,
});

export const updateLastBlockInfo = payload => ({
  type: UPDATE_LAST_BLOCK_INFO,
  payload,
});

export const connectTo = (hostname, port) => ({
  type: CONNECT,
  hostname,
  port,
});

export const connected = payload => ({
  type: CONNECTED,
  payload,
});

export const connectionError = payload => ({
  type: CONNECTION_ERROR,
  payload,
});

export const disconnected = payload => ({
  type: DISCONNECTED,
  payload,
});

// export const listenForUpdates = socket => dispatch => {
//   //dispatch(connecting());
//   const host = process.env.REACT_APP_HOST_WS || window.location.host;
//   let ws = new WebSocket(`ws:/${host}/stats`);

//   ws.onopen = () => dispatch(connected());
//   ws.onerror = () => dispatch(connectionError());
//   ws.onclose = () => dispatch(disconnected());
//   ws.onmessage = ({ data }) => {
//     const payload = JSON.parse(data);
//     const { metric, value, data: packet, timestamp } = payload;

//     switch (metric) {
//       case "cpu":
//       case "mem":
//       case "latency":
//       case "disk":
//         dispatch(updateMetrics[metric](+value, timestamp));
//         break;
//       case "log":
//         const { code, level } = packet;
//         if (code && code === "round") {
//           const { round, blockHash, blockTime } = packet;

//           const block = {
//             height: round,
//             hash: blockHash,
//             timestamp,
//           };
//           dispatch(updateLastBlockInfo(block));
//           dispatch(updateBlockTimeRead(blockTime, timestamp));
//           break;
//         }

//         if (level) {
//           const { time = timestamp } = packet;
//           dispatch(updateWarningList(packet, time));
//           break;
//         }

//       case "tail":
//         dispatch(updateLogRead(value, timestamp));
//         break;

//       default:
//         console.log("INVALID METRIC SENT");
//     }
//   };

//   if (!("ws" in socket)) {
//     socket.ws = ws;
//   }
// };
