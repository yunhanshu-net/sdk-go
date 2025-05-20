package api

//func newParamInfo(tag *tagx.FieldInfo, renderType string) (*ParamInfo, error) {
//	if tag == nil {
//		return nil, fmt.Errorf("tag==nil")
//	}
//	if tag.Tags == nil {
//		tag.Tags = map[string]string{}
//	}
//	widgetIns, err := widget.NewWidget(tag, renderType)
//	if err != nil {
//		return nil, err
//	}
//	valueType, err := tag.GetValueType()
//	if err != nil {
//		return nil, err
//	}
//	if !types.IsValueType(valueType) {
//		return nil, fmt.Errorf("不是合法的值类型：%s", valueType)
//	}
//	validate := tag.Type.Tag.Get("validate")
//	split := strings.Split(validate, ",")
//
//	param := &ParamInfo{
//		Code:         tag.Tags["code"],
//		Name:         tag.Tags["name"],
//		Desc:         tag.Tags["desc"],
//		Required:     slicesx.ContainsString(split, "required"),
//		Validates:    strings.Join(slicesx.RemoveString(split, "required"), ","),
//		Callbacks:    tag.Tags["callback"],
//		WidgetConfig: widgetIns,
//		WidgetType:   widgetIns.GetWidgetType(),
//		ValueType:    types.UseValueType(tag.Tags["type"], valueType),
//		Example:      tag.Tags["example"],
//	}
//
//	if param.Code == "" {
//		get := tag.Type.Tag.Get("json")
//		if get != "" {
//			param.Code = strings.Split(get, ",")[0]
//		}
//	}
//
//	if param.Name == "" {
//		param.Name = param.Code
//	}
//
//	return param, nil
//}
//
//func newParams(fields []*tagx.FieldInfo, renderType string) (*Params, error) {
//	//	判断不同数据类型form,table,echarts,bi,3D ....
//	children := make([]*ParamInfo, 0, len(fields))
//	for _, field := range fields {
//		info, err := newParamInfo(field, renderType)
//		if err != nil {
//			return nil, err
//		}
//		children = append(children, info)
//	}
//
//	return &Params{RenderType: stringsx.DefaultString(renderType, response.RenderTypeForm), Children: children}, nil
//}
