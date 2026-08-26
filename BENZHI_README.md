基于 Go 实现的城市轨道交通再生制动能量回收系统项目，一款城轨供电控制服务，完成再生制动能量吸收、逆变回馈、接触网电压控制与多车协同管理。

# RegenerativeBrake

RegenerativeBrake 是城市轨道交通供电控制服务。列车制动时再生能量回馈接触网，吸收装置按电压投入吸收或逆变回馈；多车同时制动时按列车状态协同分配吸收容量；接触网电压越限触发吸收，电压回落进入恢复；吸收装置故障、恢复与运行记录统一纳入告警与台账管理。

## 构建

```bash
go build -mod=vendor ./...
```

## 运行

```bash
go run -mod=vendor ./cmd/regenbrake -addr 127.0.0.1:8090 -dir ./data
```

启动后访问 http://127.0.0.1:8090/ 打开控制台页面。

## HTTP 接口

- `GET /healthz` 健康检查
- `GET /api/v1/status` 查询运行状态
- `POST /api/v1/voltage/sample` 上报接触网电压
- `POST /api/v1/trains/register` 登记列车
- `POST /api/v1/trains/braking` 更新列车制动状态
- `POST /api/v1/coop/allocate` 多车协同分配
- `POST /api/v1/devices/restore` 恢复吸收装置

## 状态机

- 吸收装置：idle -> absorbing -> discharging -> idle
- 电压控制：normal -> high -> absorbing -> restoring
- 多车协同：planning -> allocating -> settled
