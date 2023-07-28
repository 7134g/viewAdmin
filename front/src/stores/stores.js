// import { createPinia } from 'pinia';
//
// const pinia = createPinia();
//
// export const useStore = pinia.createStore({
//   state: () => ({
//     _type: '',
//   }),
//   actions: {
//     setDbType(dbType) {
//       this.db_type = dbType;
//     },
//   },
// });



import { defineStore } from 'pinia'

export const useCounterStore = defineStore('counter', {
  state: () => {
    return {
      db_type: "mysql",// 数据库
      table_name: "", // 表名
      data_struct: {}, // 表结构
      dataPage: {
        table_name: "",
        db_type: "",
        page: 1,
        size: 10,
      },
    }
  },
  // 也可以定义为
  // state: () => ({ count: 0 })
  actions: {
    setDbType(type) {
      this.db_type = type;
    },
    getDbType(){
      return this.db_type
    },


    setTableName(table_name){
      this.table_name = table_name
    },
    getTableName(){
      return this.table_name
    },

    setDataStruct(data_struct) {
      this.data_struct = data_struct
    },
    getDataStruct() {
      return this.data_struct
    },

    setDataPage(page) {
      page.db_type = this.db_type
      this.dataPage = page
    },
    getDataPage() {
      return this.dataPage
    },

  },
})
