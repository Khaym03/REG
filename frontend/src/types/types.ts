export interface DateRange {
  from: Date
  to: Date
}

export interface WorkflowInput {
  dateRange: DateRange
  receive_guides_in_transit: boolean
}
