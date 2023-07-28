<template>
  <el-scrollbar >
  <div class="back-icon" @click="closeForm">
    <el-icon >
      <Back />
    </el-icon>
  </div>

  <el-form
      ref="ruleFormRef"
      :model="ruleForm"
      :rules="rules"
      label-width="120px"
      class="demo-ruleForm"
      status-icon
  >



    <el-form-item  v-for="(key, value) in ruleForm" :label="value">
      <el-input v-model="formData[value]"/>
    </el-form-item>

    <el-button @click="printLn"></el-button>

    <el-form-item>
      <el-button type="primary" @click="submitForm(formData)">Submit</el-button>
      <el-button @click="resetForm">Reset</el-button>
    </el-form-item>

  </el-form>
  </el-scrollbar>
</template>


<style>
.back-icon {
  font-size: 25px
}
</style>

<script>

import {useCounterStore} from "../../stores/stores";
export default {
  data() {
    return {
      ruleForm: {},
      rules: {},

      formData:{}
    }
  },

  methods: {
    printLn(){
      console.log(JSON.stringify(this.formData))
    },

    submitForm(formData){

    },
    resetForm(){
      this.formData = {}
    },

    dataStruct(){
      const ds = useCounterStore().getDataStruct()
      // let keys = [];
      for (const key in ds) {
        if (ds.hasOwnProperty(key)) {
          const value = ds[key]

          this.rules[key] = [
            {
              required: true,
              message: value,
              trigger: 'change',
            }
          ]

          let typeName = ""
          let message = ""


          switch (true) {
            case value.indexOf("int") !== -1:
            case value.indexOf("float") !== -1:
              typeName = "number"
              message = value + " is number"
              break
            case value.indexOf("primitive.ObjectID") !== -1:
            case value.indexOf("string") !== -1:
            case value.indexOf("text") !== -1:
            case value.indexOf("char") !== -1:
              typeName = "string"
              message = value + " is string"
          }

          this.rules[key].push({
            type: typeName,
            message: message
          })

        }
          // keys.push(key);
      }
      console.log(JSON.stringify(this.rules))

      this.ruleForm = ds
    },

    closeForm(){
      this.$emit('closeForm-flag');
    },
  },


  mounted() {
    // 在其他方法或是生命周期中也可以调用方法
    this.dataStruct()
  }
}

</script>


