package widget

import (
	"github.com/yunhanshu-net/sdk-go/pkg/stringsx"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
)

func NewWidget(info *tagx.FieldInfo) (Widget, error) {
	if info.Tags == nil {
		info.Tags = make(map[string]string)
	}
	widgetType := stringsx.DefaultString(info.Tags["widget"], "input")

	switch widgetType {
	case WidgetInput:
		return NewInputWidget(info)
		//case WidgetCheckbox:
		//	return newCheckboxWidget(info)
		//case WidgetRadio:
		//	return newRadioWidget(info)
		//case WidgetSelect:
		//	return newSelectWidget(info)
		//case WidgetSwitch:
		//	return newSwitchWidget(info)
		//case WidgetSlider:
		//	return newSliderWidget(info)
		//case WidgetFile:
		//	return newFileWidget(info)
	}
	return NewInputWidget(info)
}
