import { createSelector } from "reselect";

const nodeUpdatesSelector = state => state.updates;
const lastBlockInfo = state => state.lastBlock;
const dusk1UpdatesSelector = state => state.updates["dusk-1"] || [];
const regionSelector = state => state.regions || {};
const cpuSelector = state => state.cpu || [];
const logSelector = state => state.log || [];
const netSelector = state => state.net || [];
const memSelector = state => state.memory || [];
const diskSelector = state => state.disk || [];

const last20dusk1UpdatesSelector = createSelector(
  dusk1UpdatesSelector,
  updates => updates.slice(0, 20).reverse()
);

export const lastNodeUpdateSelector = createSelector(
  nodeUpdatesSelector,
  items =>
    Object.entries(items).map(([hostname, updates]) => ({
      hostname,
      ...updates[0]
    }))
);

export const getCurrentBlockInfo = createSelector(
  lastBlockInfo,
  info => info
);

export const getCPUMetrics = createSelector(
  cpuSelector,
  info =>
    info
      .slice(0, 20)
      .reverse()
      .reduce(
        (acc, { value, timestamp }) => {
          acc.labels.push(timestamp);
          acc.series[0].push(value);
          return acc;
        },
        { labels: [], series: [[]] }
      )
);

export const getNetMetrics = createSelector(
  netSelector,
  info =>
    info
      .slice(0, 20)
      .reverse()
      .reduce(
        (acc, { value, timestamp }) => {
          acc.labels.push(timestamp);
          acc.series[0].push(value);
          return acc;
        },
        { labels: [], series: [[]] }
      )
);

export const getLogMetrics = createSelector(
  logSelector,
  info => info.slice(0, 200)
);

export const getMemoryMetrics = createSelector(
  memSelector,
  info =>
    info
      .slice(0, 20)
      .reverse()
      .reduce(
        (acc, { value, timestamp }) => {
          acc.labels.push(timestamp);
          acc.series[0].push(value);
          return acc;
        },
        { labels: [], series: [[]] }
      )
);

// export const getDiskMetrics = createSelector(
//   diskSelector,
//   info => info
// );

export const getBlockTime = createSelector(
  last20dusk1UpdatesSelector,
  updates => updates.map(({ time, height }) => ({ time, height }))
);

export const getNetLatency = createSelector(
  last20dusk1UpdatesSelector,
  updates => updates.map(({ ping, height }) => ({ ping, height }))
);

export const getNodeLocations = createSelector(
  regionSelector,
  regions =>
    Object.entries(regions).map(([region, nodes]) => ({
      name: region,
      value: nodes.length
    }))
);

export const getHighestScore = createSelector(
  last20dusk1UpdatesSelector,
  updates => updates.map(({ score, height }) => ({ score, height }))
);

export const getDiskMetrics = createSelector(
  diskSelector,
  info => {
    if (info[0]) {
      let used = info[0].value;
      let free = 100 - used;
      return {
        series: [used, free],
        labels: [`Used: ${used}%`, `Free: ${free}%`]
      };
    }
  }
);

// const data = [
//   { name: "ams1", value: 400 },
//   { name: "ams2", value: 300 },
//   { name: "sfo", value: 300 },
//   { name: "Lon", value: 200 }
// ];
