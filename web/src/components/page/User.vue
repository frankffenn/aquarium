<template>
    <div>
        <div class="crumbs">
        </div>
        <div class="container">
            <div class="handle-box">
                <el-button
                    type="primary"
                    class="handle-del mr10"
                    size="mini"
                >Reload</el-button>
                <el-button
                    class="handle-del mr10"
                    size="mini"
                >Add</el-button>
                <el-button
                    class="handle-del mr10"
                    @click="handleDelete"
                    size="mini"
                >Delete</el-button>
            </div>
            <el-table
                :data="tableData"
                class="table"
                ref="multipleTable"
                empty-text="No Data"
                header-cell-class-name="table-header"
                @selection-change="handleSelectionChange"
                @row-click="handleEdit"
                >
                <el-table-column type="selection" width="55" align="center"></el-table-column>
                <el-table-column v-if="false" prop="ID" label="ID" width="55" align="center"></el-table-column>
                <el-table-column prop="Username" label="Username"></el-table-column>
                <el-table-column prop="Level" label="Level"></el-table-column>
                <el-table-column prop="CreatedAt" label="CreatedAt" :formatter="dateFormat"></el-table-column>
                <el-table-column prop="UpdatedAt" label="UpdatedAt" :formatter="dateFormat"></el-table-column>

            </el-table>
            <div class="pagination">
                <el-pagination
                    background
                    layout="total, prev, pager, next"
                    :current-page="query.pageIndex"
                    :page-size="query.pageSize"
                    :total="pageTotal"
                    @current-change="handlePageChange"
                ></el-pagination>
            </div>
        </div>

        <!-- 编辑弹出框 -->
        <el-dialog title="User-" :visible.sync="editVisible" width="30%">
            <el-form ref="form" :model="form" label-width="100px" label-position="left">
                <el-form-item label="Username:">
                    <el-input v-model="form.Username"></el-input>
                </el-form-item>
                <el-form-item label="Level:">
                    <el-input v-model="form.Level"></el-input>
                </el-form-item>
                 <!-- <el-form-item label="Password">
                    <el-input v-model="form.Password"></el-input>
                </el-form-item>
                 <el-form-item label="Repeat">
                    <el-input v-model="form.Repeat"></el-input>
                </el-form-item> -->
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="editVisible = false">Cancel</el-button>
                <el-button type="primary" @click="saveEdit">OK</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script>
import { formatDateTime } from "../../utils/date";
import { userListReq } from '../../api/index';
export default {
    name: 'basetable',
    data() {
        return {
            query: {
                page: 1,
                size: 10
            },
            tableData: [],
            multipleSelection: [],
            delList: [],
            editVisible: false,
            pageTotal: 0,
            form: {},
            idx: -1,
            id: -1
        };
    },
    created() {
        this.getUserList();
    },
    methods: {
         dateFormat(row, column, cellValue, index){
            if (!cellValue) {
                return ''
            }
            var date = new Date(cellValue);
            return formatDateTime(date);
        }, 
        nameStyle(column){
            if (column.columnIndex === 2) {
                return 'color:#108EE9';
            }
        },             
        getUserList() {
            var token = localStorage.getItem("token");
            userListReq(this.query,token).then(res => {
                if (res.success) {
                    this.tableData = res.data.users;
                    this.pageTotal = res.data.total || 50;
                }else{
                    if (res.code == 401) {   
                        localStorage.removeItem('ms_username');
                        this.$router.push("/login")
                    }
                    this.$message.error(res.msg || "unkown err");
                }
            });
        },
        // 触发搜索按钮
        handleSearch() {
            this.$set(this.query, 'pageIndex', 1);
            this.getData();
        },
        // 删除操作
        handleDelete(index, row) {
            // 二次确认删除
            this.$confirm('Are you sure to DELETE ? ', '', {
                type: 'warning'
            }).then(() => {
                   this.delAllSelection();
            }).catch(() => {});
        },
        // 多选操作
        handleSelectionChange(val) {
            this.multipleSelection = val;
        },
        delAllSelection() {
            const length = this.multipleSelection.length;
            if (length <= 0) {
                return
            }
            let str = '';
            this.delList = this.delList.concat(this.multipleSelection);
            for (let i = 0; i < length; i++) {
                str += this.multipleSelection[i].name + ' ';
            }
            this.$message.error(`Deleted ${str}`);
            this.multipleSelection = [];
        },
        // 编辑操作
        handleEdit(row,column,event) {
            this.form = row;
            this.editVisible = true;
        },
        // 保存编辑
        saveEdit() {
            this.editVisible = false;
            this.$message.success(`Updated ${this.idx + 1} line done`);
            this.$set(this.tableData, this.idx, this.form);
        },
        // 分页导航
        handlePageChange(val) {
            this.$set(this.query, 'pageIndex', val);
            this.getData();
        }
    }
};
</script>

<style scoped>
.handle-box {
    margin: 10px 10px;
}

.handle-select {
    width: 120px;
}

.handle-input {
    width: 300px;
    display: inline-block;
}
.table {
    width: 100%;
    font-size: 12px;
}
.red {
    color: #ff0000;
}
.mr10 {
    margin-right: 10px;
}
.table-td-thumb {
    display: block;
    margin: auto;
    width: 40px;
    height: 40px;
}
</style>
