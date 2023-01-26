import { JsonObject, JsonProperty } from 'typescript-json-serializer';

@JsonObject()
export class Property {
  @JsonProperty() ID: number = 0
  @JsonProperty() CreatedAt: Date = new Date()
  @JsonProperty() UpdatedAt: Date = new Date()
  @JsonProperty() DeletedAt: Date | undefined = new Date()
  @JsonProperty() Url: string = ""
  @JsonProperty() Price: number = 0
  @JsonProperty() LandArea: number = 0
  @JsonProperty() BuildingArea: number = 0
  @JsonProperty() Station: string = ""
  @JsonProperty() City: string = ""
  @JsonProperty() Layout: string = ""
  @JsonProperty() BuildYear: number = 0
  @JsonProperty() Access: string = ""
  @JsonProperty() Road: string = ""
  @JsonProperty() OtherCost: string = ""
  @JsonProperty() CoverageRatio: string = ""
  @JsonProperty() Timing: string = ""
  @JsonProperty() Rights: string = ""
  @JsonProperty() Structure: string = ""
  @JsonProperty() BuildCompany: string = ""
  @JsonProperty() Reform: string = ""
  @JsonProperty() LandKind: string = ""
  @JsonProperty() AreaPurpose: string = ""
  @JsonProperty() OtherRestriction: string = ""
  @JsonProperty() OtherNotice: string = ""
  @JsonProperty() JobId: number = 0
  @JsonProperty() AreaPrice: number = 0
  @JsonProperty() EstimatedBuildingPrice: number = 0
  @JsonProperty() ClickCount: number = 0
}

@JsonObject()
export class ResponseProperties {
  @JsonProperty() count: number = 0
  @JsonProperty() properties: Property[] = []
}