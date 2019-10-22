/**
 * Created by solee on 1/26/16.
 */
'use strict';

var helper = require('../../utils/helper');

module.exports = function (sequelize, DataTypes) {
  return sequelize.define('package_activation_records', {
    // 激活记录id
    id: {type: DataTypes.INTEGER(11), allowNull: false, autoIncrement: true, primaryKey: true, field: 'id'},
    // 用户id
    user_id: {type: DataTypes.INTEGER(11), allowNull: false, field: 'user_id'},
    // 服务礼包id
    service_package_id: {type: DataTypes.INTEGER(11), allowNull: false, field: 'service_package_id'},
    // 激活渠道 0:未知 1:后台发放 2:用户自开通 3:购买激活
    channel: {type: DataTypes.INTEGER(1), allowNull: true, field: 'channel'},
    // 礼包内容备份
    backup: {type: DataTypes.TEXT, allowNull: true, field: 'backup'},
    // 礼包价格备份
    backup_price: {type: DataTypes.INTEGER(11), allowNull: true, field: 'backup_price'},
    // 来源管理员id
    from_operator_id: {type: DataTypes.INTEGER(11), allowNull: true, field: 'from_operator_id'},
    // 是否是vip 0:否（默认）1:是
    is_vip: {type: DataTypes.INTEGER(11), allowNull: true, field: 'is_vip'},
    // 激活状态 1:未激活 2:激活中 3:已激活
    status: {type: DataTypes.INTEGER(1), allowNull: true, field: 'status'},
    // 车位id
    carport_id: {type: DataTypes.INTEGER(11), allowNull: true, field: 'carport_id'},
    // 消费订单id
    consume_order_id: {type: DataTypes.INTEGER(11), allowNull: true, field: 'consume_order_id'},
    // 发放记录id
    package_allocate_record_id: {type: DataTypes.INTEGER(11), allowNull: true, field: 'package_allocate_record_id'},
    // 备注
    remark: {type: DataTypes.STRING(50), allowNull: true, field: 'remark'},
    // 创建时间
    created_at: {type: DataTypes.INTEGER(11), allowNull: true, field: 'created_at'},
    // 更新时间
    updated_at: {type: DataTypes.INTEGER(11), allowNull: true, field: 'updated_at'}
  }, {
    timestamps: false,
    freezeTableName: true,
    hooks: {
      beforeCreate: function (instance, options) {
        instance.created_at = instance.updated_at = helper.currentTimestamp();
      },
      beforeUpdate: function (instance, options) {
        instance.updated_at = helper.currentTimestamp();
      },
      beforeBulkUpdate: function (options) {
        options.individualHooks = true;
      },
      beforeBulkCreate: function (instances, options) {
        options.individualHooks = true;
      }
    }
  });
};