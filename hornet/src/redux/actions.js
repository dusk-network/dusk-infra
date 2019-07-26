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
  UPDATE_THREAD,
  UPDATE_TIME_READ,
  UPDATE_TX_READ,
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

export const updateTxNr = (value, timestamp) => ({
  type: UPDATE_TX_READ,
  value,
  timestamp,
});

export const updateWarningList = (value, timestamp) => ({
  type: UPDATE_WARN_LIST,
  value,
  timestamp,
});

export const updateThread = (value, timestamp) => ({
  type: UPDATE_THREAD,
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

export const connectTo = (hostname, port = "") => ({
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
