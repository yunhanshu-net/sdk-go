package render

import (
	"fmt"
	"github.com/yunhanshu-net/sdk-go/pkg/tagx"
)

func NewWidget(info *tagx.FieldInfo) (Widget, error) {
	if info.Tags == nil {
		return nil, fmt.Errorf("tags==nil")
	}
	widgetType := "input"
	if info.Tags["widget"] != "" {
		widgetType = info.Tags["widget"]
	}

	switch widgetType {
	case WidgetInput:
		return newInputWidget(info)
	case WidgetCheckbox:
		return newCheckboxWidget(info)
	case WidgetRadio:
		return newRadioWidget(info)
	case WidgetSelect:
		return newSelectWidget(info)
	case WidgetSwitch:
		return newSwitchWidget(info)
	case WidgetSlider:
		return newSliderWidget(info)
	case WidgetFile:
		return newFileWidget(info)
	}
	return newInputWidget(info)
}
