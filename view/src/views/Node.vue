<template>
  <section>
    <h1>Bastion</h1>
    <table>
      <thead>
        <tr>
          <th>Hostname</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-bind:key="bastion" v-for="bastion in bastions">
          <td>{{ bastion }}</td>
          <td><a href="javascript:" @click="deleteInstance(bastion)">delete</a></td>
        </tr>
      </tbody>
    </table>

    <br>
    <input type="text" v-model="bastionName">
    <button @click="newInstance(true)">add</button>

    <h1>Nodes</h1>
    <table>
      <thead>
        <tr>
          <th>Hostname</th>
          <th>Power</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-bind:key="node" v-for="node in nodes">
          <td>{{ node }}</td>
          <td>
            <a href="javascript:" @click="powerControl(node, 'stop')">stop</a>,
            <a href="javascript:" @click="powerControl(node, 'restart')">restart</a>
          </td>
          <td>
            <a href="javascript:" @click="deleteInstance(node)">delete</a>,
            <a href="javascript:" @click="openConfig(node)">config</a>
          </td>
        </tr>
      </tbody>
    </table>

    <br>
    <input type="text" v-model="nodeName">
    <button @click="newInstance(false)">add</button>

    <div v-show="settingInstance">
      <h1>Config</h1>
      <textarea v-model="config"></textarea><br>
      <button @click="applyConfig()">apply</button>
    </div>
  </section>
</template>

<style>
section {
  padding: 1em 2em;
}
table {
  border-collapse: collapse;
}
th, td {
  padding: .5em 2em;
  border: 1px solid #999;
}
textarea {
  width: 70vw;
  height: 25em;
}
</style>

<script>
"use strict";

export default {
  data() {
    return {
      bastions: [],
      nodes: [],

      bastionName: "",
      nodeName: "",

      settingInstance: null,
      config: "",
    }
  },
  created() {
    this.loadInstances();
  },
  methods: {
    async loadInstances() {
      const bResp = await fetch(`/api/database/query`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify("SELECT host FROM instances WHERE bastion = 1"),
      });
      if(bResp.status != 200){
        return alert(await bResp.text());
      }
      this.bastions = (await bResp.json()).map(({host}) => host);

      const resp = await fetch(`/api/collector/instance`);
      if(resp.status != 200){
        return alert(await resp.text());
      }
      this.nodes = await resp.json();
    },
    async newInstance(bastion) {
      const name = bastion ? this.bastionName : this.nodeName;
      if(!name){
        return alert("name is null");
      }

      const resp = await fetch(`/api/collector/instance/${name}${bastion ? "?bastion=1" : ""}`, {method: "PUT"});
      if(resp.status != 200){
        return alert(await resp.text());
      }

      return this.loadInstances();
    },
    async deleteInstance(node) {
      if(!confirm(`delete ${node}`)){
        return;
      }

      const resp = await fetch(`/api/collector/instance/${node}`, {method: "DELETE"});
      if(resp.status != 200){
        return alert(await resp.text());
      }

      return this.loadInstances();
    },
    async openConfig(node) {
      const resp = await fetch(`/api/cnc/${node}/state`);
      if(resp.status != 200){
        return alert(await resp.text());
      }

      this.config = JSON.stringify(await resp.json(), null, 2);
      this.settingInstance = node;
    },
    async applyConfig() {
      try {
        const parsed = JSON.parse(this.config);

        const resp = await fetch(`/api/cnc/${this.settingInstance}/state`, {
          method: "PUT",
          headers: {"Content-Type": "application/json"},
          body: JSON.stringify(parsed),
        });
        if(resp.status != 200){
          throw await resp.text();
        }

        this.config = "";
        this.settingInstance = null;
      } catch(e) {
        return alert(e);
      }
    },
    async powerControl(node, mode) {
      if(!confirm(`${mode} ${node}`)){
        return;
      }

      const resp = await fetch(`/api/cnc/${node}/power`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(mode),
      });
      if(resp.status != 200){
        throw await resp.text();
      }
    }
  }
}
</script>
