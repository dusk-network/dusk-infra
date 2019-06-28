/* eslint-disable no-script-url */

import React from "react";
import Link from "@material-ui/core/Link";
import { makeStyles } from "@material-ui/styles";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Title from "./Title";

const useStyles = makeStyles(theme => ({}));

export default function Nodes({ items = [] }) {
  const classes = useStyles();
  return (
    <>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>Hostname</TableCell>
            <TableCell>IP</TableCell>
            <TableCell>Region</TableCell>
            <TableCell>Height</TableCell>
            <TableCell>Stake</TableCell>
            <TableCell>Hash</TableCell>
            <TableCell>Ping</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {items.map(update => (
            <TableRow key={update.hostname}>
              <TableCell>{update.hostname}</TableCell>
              <TableCell>{update.IP}</TableCell>
              <TableCell>{update.region}</TableCell>
              <TableCell>{update.height}</TableCell>
              <TableCell>{update.stake}</TableCell>
              <TableCell>{update.hash}</TableCell>
              <TableCell>{update.ping}ms</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </>
  );
}
